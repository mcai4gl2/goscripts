package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/mcai4gl2/goscripts/commoncrawl"
)

func getStockProfileForAllTickers(tickerFileName string, outputFileName string, crawlParallel int) {
	saveChan := make(chan CrawlResult, 5)
	rand.Seed(time.Now().UnixNano())

	var crawlWaitGroup sync.WaitGroup
	crawlWaitGroup.Add(crawlParallel)

	tickerChan := make(chan string)
	go func() {
		tickers := crawl.GetAllTickers(tickerFileName, "ticker")
		for ticker := range tickers {
			tickerChan <- ticker
		}
		close(tickerChan)
	}()

	for i := 0; i < crawlParallel; i++ {
		go func(index int) {
			defer crawlWaitGroup.Done()
			for ticker := range tickerChan {
				log.Println(fmt.Sprintf("Start processing for ticker %s", ticker))
				url := formatUrl(ticker)
				data, err := getData(ticker, url)
				result := extractData(ticker, data, err)
				saveChan <- result
				time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
			}
		}(i)
	}

	var saveWaitGroup sync.WaitGroup
	saveWaitGroup.Add(1)

	go func() {
		defer saveWaitGroup.Done()

		file, _ := os.Create(outputFileName)
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"ticker", "name", "sector", "industry"}
		writer.Write(header)

		for result := range saveChan {
			log.Println(fmt.Sprintf("Ticker: %s, Name: %s, Sector: %s, Industry: %s",
				result.ticker, result.data["name"],
				result.data["sector"], result.data["industry"]))
			writer.Write([]string{result.ticker,
				result.data["name"],
				result.data["sector"],
				result.data["industry"]})
		}
	}()

	log.Println("Waiting for crawl tasks to finish")
	crawlWaitGroup.Wait()
	close(saveChan)
	log.Println("Waiting for save tasks to finish")
	saveWaitGroup.Wait()
}

func main() {
	tickerFilePtr := flag.String("ticker", "", "Ticker input file")
	outputFilePrt := flag.String("output", "", "Full path to output file")
	webParallelPtr := flag.Int("webParallel", 5, "Number of concurrent go routine to crawl from yahoo")

	flag.Parse()

	if *tickerFilePtr == "" {
		panic("ticker file cannot be empty")
	}

	if *outputFilePrt == "" {
		panic("output file cannot be empty")
	}

	getStockProfileForAllTickers(*tickerFilePtr, *outputFilePrt, *webParallelPtr)
}
