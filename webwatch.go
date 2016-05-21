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

func getTextFromPage(url, cssSelector string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	arr := doc.Find(cssSelector).Map(func(i int, s *goquery.Selection) string {
		return strings.TrimSpace(s.Text())
	})

	fmt.Printf("Hi %v", arr)
}

func loadUrls(filePath string) {
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
        fmt.Println(rules[len(rules)-1], line)
        rules[len(rules)-1].filter = line
      }
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

  fmt.Println(rules)
}

func main() {
	pathToUrls := flag.String("config", "urls.txt", "path to file with urls")
	loadUrls(*pathToUrls)
	// getTextFromPage("http://apple.multitronic.fi/en/products/1398619/iphone-se-64gb-space-grey", "#prod_stock span")
	// getTextFromPage("http://apple.multitronic.fi/en/products/1398619/iphone-se-64gb-space-grey", ".mt_well_light label")
}
