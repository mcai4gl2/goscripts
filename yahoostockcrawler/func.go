package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

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

func getUrl(url string) (data string, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprintf("Failed process url: %s with error %s", url, r))
			data = ""
			err = r.(error)
		}
	}()

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	response, _ := netClient.Get(url)

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	bodyStr := string(body)

	bodyStr = bodyStr[strings.Index(bodyStr, "HistoricalPriceStore"):]

	bodyStr = bodyStr[0 : strings.Index(bodyStr, "]")+1]

	data = bodyStr[strings.Index(bodyStr, "["):]
	err = nil

	return
}

func getFullOutputFileName(ticker string, startDate string, endDate string,
	outputDir string) string {
	return path.Join(outputDir,
		fmt.Sprintf("%s_%s_%s.json", ticker, startDate, endDate))
}

func writeToFile(fullOutputFileName string, content string) error {
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
