package main

import (
	"fmt"

	"github.com/anaskhan96/soup"
)

func formatUrl(ticker string) string {
	return fmt.Sprintf("https://finance.yahoo.com/quote/%s/profile?p=%s",
		ticker, ticker)
}

type CrawlResult struct {
	ticker   string
	sector   string
	industry string
	err      error
}

func extractData(ticker string, web string, err error) CrawlResult {
	if err != nil {
		return CrawlResult{ticker, "", "", err}
	}

	doc := soup.HTMLParse(web)

	nodes := doc.FindAllStrict("span")
	sector := ""
	for _, node := range nodes {
		if node.Text() == "Sector" {
			sector = node.FindNextElementSibling().Text()
		}
	}
	industry := ""
	for _, node := range nodes {
		if node.Text() == "Industry" {
			industry = node.FindNextElementSibling().Text()
		}
	}
	return CrawlResult{ticker, sector, industry, nil}
}

func getData(ticker string, url string) (string, error) {
	resp, err := soup.Get(url)
	return resp, err
}
