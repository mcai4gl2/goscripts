package main

import "testing"
import "fmt"

func TestParseDateString(t *testing.T) {
	time, err := parseDateString("20190827")
	if err != nil {
		t.Errorf("There shall be no error, but we got %s", err)
	}
	if time.Unix() != 1566864000 {
		t.Errorf("Expecting the time to be matching 1566864000")
	}
}

func TestFormatUrl(t *testing.T) {
	url, _ := formatUrl("6862.HK", "20190101", "20190827")
	fmt.Println(url)
}
