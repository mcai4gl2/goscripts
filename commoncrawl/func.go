package crawl

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/mcai4gl2/goscripts/util"
)

func GetAllTickers(tickerFileFullPath string, tickerColumnName string) <-chan string {
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
