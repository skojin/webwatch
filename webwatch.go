package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
)

type WebsiteRule struct {
	url    string
	filter string
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
			rules = append(rules, WebsiteRule{url: scanner.Text()})
		} else {
			if len(rules) > 0 && len(line) > 0 {
				rules[len(rules)-1].filter = line
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return rules
}

func checkEachWebsite(rules []WebsiteRule, testMode bool) {
	for _, rule := range rules {
		t := getTextFromPage(rule.url, rule.filter)
		if testMode {
			fmt.Printf("%s\n   %s\n\n", rule.url, t)
		} else {

		}
	}
}

func main() {
	pathToUrlsPtr := flag.String("config", "urls.txt", "path to file with urls")
	testModePtr := flag.Bool("test", false, "check each site and output result")
	flag.Parse()

	checkEachWebsite(loadUrlRules(*pathToUrlsPtr), *testModePtr)
	// fmt.Printf("%v", rules)
	// v := getTextFromPage("http://apple.multitronic.fi/en/products/1398619/iphone-se-64gb-space-grey", "#prod_stock span")
	// fmt.Println(v)
	// getTextFromPage("http://apple.multitronic.fi/en/products/1398619/iphone-se-64gb-space-grey", ".mt_well_light label")
}
