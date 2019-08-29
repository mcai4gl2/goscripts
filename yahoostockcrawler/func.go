package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mcai4gl2/goscripts/util"
)

func getAllHKEXTickers(tickerFileFullPath string, tickerColumnName string) <-chan string {
	tickers := make(chan string)

	go func() {
		tickerFile, err := os.Open(tickerFileFullPath)
		defer tickerFile.Close()

		defer close(tickers)

		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(tickerFile)

		r := csv.NewReader(reader)

		firstLine, _ := r.Read()
		indexOfTicker := util.FirstIndexOf(firstLine, tickerColumnName)

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			tickers <- record[indexOfTicker]
		}
	}()

	return tickers
}

func parseDateString(date string) (time.Time, error) {
	layout := "20060102"
	return time.Parse(layout, date)
}

func formatUrl(ticker string, startDate string, endDate string) (string, error) {
	start, err := parseDateString(startDate)
	if err != nil {
		return "", err
	}
	end, err := parseDateString(endDate)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://finance.yahoo.com/quote/%s/history?period1=%d&period2=%d&interval=1d&filter=history&frequency=1d",
		ticker, start.Unix(), end.Unix()), nil
}

func getUrl(url string) string {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	response, _ := netClient.Get(url)

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	bodyStr := string(body)

	bodyStr = bodyStr[strings.Index(bodyStr, "HistoricalPriceStore"):]

	bodyStr = bodyStr[0 : strings.Index(bodyStr, "]")+1]
	bodyStr = bodyStr[strings.Index(bodyStr, "["):]

	return bodyStr
}

func getFullOutputFileName(ticker string, startDate string, endDate string,
	outputDir string) string {
	return path.Join(outputDir,
		fmt.Sprintf("%s_%s_%s.json", ticker, startDate, endDate))
}

func writeToFile(ticker string, startDate string, endDate string,
	outputDir string, content string) error {
	fullOutputFileName := getFullOutputFileName(ticker, startDate,
		endDate, outputDir)

	file, err := os.Create(fullOutputFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}
	return file.Sync()
}
