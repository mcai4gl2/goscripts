package main

import (
	"fmt"
)

func main() {
	tickers := getAllHKEXTickers("/home/ligeng/Codes/scripts/sec_list.csv", "ticker")

	for ticker := range tickers {
		fmt.Println(ticker)
	}
}
