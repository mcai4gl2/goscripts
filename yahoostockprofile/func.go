package main

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

func formatUrl(ticker string) string {
	return fmt.Sprintf("https://finance.yahoo.com/quote/%s/profile?p=%s",
		ticker, ticker)
}

type CrawlResult struct {
	ticker string
	data   map[string]string
	err    error
}

func extractData(ticker string, web string, err error) CrawlResult {
	if err != nil {
		return CrawlResult{ticker, nil, err}
	}

	doc := soup.HTMLParse(web)

	data := make(map[string]string)

	nodes := doc.FindAllStrict("span")
	data["sector"] = ""
	for _, node := range nodes {
		if node.Text() == "Sector" {
			data["sector"] = node.FindNextElementSibling().Text()
		}
	}
	data["industry"] = ""
	for _, node := range nodes {
		if node.Text() == "Industry" {
			data["industry"] = node.FindNextElementSibling().Text()
		}
	}
	headers := doc.FindAllStrict("h1")
	data["name"] = ""
	for _, node := range headers {
		if strings.Contains(node.Text(), ticker) {
			data["name"] = node.Text()
		}
	}
	return CrawlResult{ticker, data, nil}
}

func getData(ticker string, url string) (string, error) {
	resp, err := soup.Get(url)
	return resp, err
}
