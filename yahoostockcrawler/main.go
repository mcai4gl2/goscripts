package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type CrawlWork struct {
	ticker         string
	startDate      string
	endDate        string
	outputFileName string
}

type SaveWork struct {
	outputFileName string
	data           string
}

func getPricesForAllTickers(tickerFile string, startDate string, endDate string,
	ouputDir string, crawlParallel int, saveParallel int, oneTickerFilter string) {
	saveChan := make(chan SaveWork)
	rand.Seed(time.Now().UnixNano())

	tickerChan := func() <-chan CrawlWork {
		tickerChan := make(chan CrawlWork)
		go func() {
			tickers := getAllHKEXTickers(tickerFile, "ticker")
			for ticker := range tickers {
				if oneTickerFilter != "" && ticker != oneTickerFilter {
					continue
				}
				outputFileName := getFullOutputFileName(ticker, startDate, endDate, ouputDir)
				if _, err := os.Stat(outputFileName); os.IsNotExist(err) {
					log.Println(fmt.Sprintf("Crawling ticker: %s", ticker))
					tickerChan <- CrawlWork{ticker, startDate, endDate, outputFileName}
				}
			}
			close(tickerChan)
		}()
		return tickerChan
	}()

	crawlOutputChannels := make([]<-chan SaveWork, crawlParallel)
	for i := 0; i < crawlParallel; i++ {
		crawlOutputChannels[i] = func() <-chan SaveWork {
			outputChan := make(chan SaveWork)
			go func() {
				for work := range tickerChan {
					url, _ := formatUrl(work.ticker, work.startDate, work.endDate)
					log.Println(fmt.Sprintf("Calling url: %s", url))
					data, err := getUrl(url)
					if err != nil {
						log.Println(fmt.Sprintf("Failing to get data for ticker %s with error %s",
							work.ticker, err))
					} else {
						log.Println("Got data, scheduling save work")
						outputChan <- SaveWork{work.outputFileName, data}
					}
					time.Sleep(time.Duration(rand.Intn(6)) * time.Second)
				}
				close(outputChan)
			}()
			return outputChan
		}()
	}

	log.Println("Creating fan-in channels")
	var crawlWaitGroup sync.WaitGroup
	crawlWaitGroup.Add(len(crawlOutputChannels))
	for _, ch := range crawlOutputChannels {
		go func(channel <-chan SaveWork) {
			defer crawlWaitGroup.Done()
			for data := range channel {
				log.Println("Fan-in to main save channel")
				saveChan <- data
			}
		}(ch)
	}

	var saveWaitGroup sync.WaitGroup
	saveWaitGroup.Add(saveParallel)
	for i := 1; i <= saveParallel; i++ {
		go func() {
			defer saveWaitGroup.Done()
			for saveWork := range saveChan {
				log.Println(fmt.Sprintf("Saving result to file: %s", saveWork.outputFileName))
				err := writeToFile(saveWork.outputFileName, saveWork.data)
				if err != nil {
					log.Println(fmt.Sprintf("Failing to save results to file: %s", err))
				}
			}
		}()
	}

	log.Println("Waiting for crawl tasks to finish")
	crawlWaitGroup.Wait()
	close(saveChan)
	log.Println("Waiting for save tasks to finish")
	saveWaitGroup.Wait()
}

func main() {
	tickerFilePtr := flag.String("ticker", "", "Ticker input file")
	startDatePtr := flag.String("start", "", "Start date in YYYYMMDD format")
	endDatePtr := flag.String("end", "", "End date in YYYYMMDD format")
	outputDir := flag.String("output", "", "Result output directory")
	webParallelPtr := flag.Int("webParallel", 10, "Number of concurrent go routine to crawl from yahoo")
	diskParallelPtr := flag.Int("diskParallel", 4, "Number of concurrent go routine to save results to disk")
	oneTickerPtr := flag.String("filter", "", "One ticker to filter on")

	flag.Parse()

	if *tickerFilePtr == "" {
		panic("ticker file cannot be empty")
	}

	if *startDatePtr == "" {
		panic("start date cannot be empty")
	}

	if *endDatePtr == "" {
		panic("end date cannot be empty")
	}

	if *outputDir == "" {
		panic("output directory cannot be empty")
	}

	getPricesForAllTickers(*tickerFilePtr, *startDatePtr, *endDatePtr,
		*outputDir, *webParallelPtr, *diskParallelPtr, *oneTickerPtr)
}
