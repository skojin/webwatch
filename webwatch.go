package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type WebsiteRule struct {
	Url    string
	Filter string
}
type WebsiteValue struct {
	Url   string
	Value string
}

// @param filter css selector or "> shell command"
func getMatchedTextFromPage(url, filter string) string {
	if strings.HasPrefix(filter, ">") {
		command := strings.TrimSpace(filter[1:])
		return executeCommand(command, downloadPage(url))
	} else {
		return getMatchedCssTextFromPage(url, filter)
	}
}

func downloadPage(url string) string {
	// get and convert to utf8
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(utf8)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func executeCommand(command, pageBody string) string {
	cmd := exec.Command("bash", "-c", command)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, pageBody)
	}()

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}

func getMatchedCssTextFromPage(url, cssSelector string) string {
	// get and convert to utf8
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	// parse
	doc, err := goquery.NewDocumentFromReader(utf8)
	if err != nil {
		panic(err)
	}
	// css query
	arr := doc.Find(cssSelector).Map(func(i int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	})
	return strings.Join(arr, "\n")
}

func loadUrlRules(filePath string) []WebsiteRule {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var rules []WebsiteRule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "http") {
			rules = append(rules, WebsiteRule{Url: scanner.Text()})
		} else if strings.HasPrefix(line, " ") {
			// comment line starts with space
		} else {
			// if we already have rules and current line not empty, then set this filter to last rule
			if len(rules) > 0 && len(line) > 0 {
				rules[len(rules)-1].Filter = line
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return rules
}

func checkEachWebsite(rules []WebsiteRule, dbPath string, testMode bool) {
	newvalues := []WebsiteValue{}
	db := loadValueDb(dbPath)
	for _, rule := range rules {
		t := getMatchedTextFromPage(rule.Url, rule.Filter)
		t_hash := fmt.Sprintf("%x", sha1.Sum([]byte(t)))
		if testMode {
			fmt.Printf("%s\n%q\n\n", rule.Url, t)
		} else {
			if t_hash != db[rule.Url] {
				fmt.Printf("%s has been updated\n", rule.Url)
			}
			newvalues = append(newvalues, WebsiteValue{rule.Url, t_hash})
		}
	}
	updateValueDb(dbPath, newvalues)
}

func loadValueDb(dbPath string) map[string]string {
	dat, err := ioutil.ReadFile(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}
		} else {
			panic(err)
		}
	}

	var values []WebsiteValue
	json.Unmarshal(dat, &values)

	db := map[string]string{}
	for _, v := range values {
		db[v.Url] = v.Value
	}
	return db
}
func updateValueDb(dbPath string, values []WebsiteValue) {
	b, err := json.Marshal(values)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(dbPath, b, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	pathToUrlsPtr := flag.String("config", "urls.txt", "path to file with urls")
	pathToDbPtr := flag.String("db", "webwatch.db", "path to database file")
	testModePtr := flag.Bool("test", false, "check each site and output result")
	flag.Parse()

	checkEachWebsite(loadUrlRules(*pathToUrlsPtr), *pathToDbPtr, *testModePtr)
}
