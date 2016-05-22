package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"os"
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

func getTextFromPage(url, cssSelector string) string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

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
		} else {
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
		t := getTextFromPage(rule.Url, rule.Filter)
		if testMode {
			fmt.Printf("%s\n   %s\n\n", rule.Url, t)
		} else {
			if t != db[rule.Url] {
				fmt.Printf("%s has beed updated\n", rule.Url)
			}
			newvalues = append(newvalues, WebsiteValue{rule.Url, t})
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
