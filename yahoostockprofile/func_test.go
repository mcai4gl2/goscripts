package main

import "testing"
import "errors"
import "github.com/stretchr/testify/assert"

func TestFormatUrl(t *testing.T) {
	url := formatUrl("0382.HK")
	assert.Equal(t, "https://finance.yahoo.com/quote/0382.HK/profile?p=0382.HK", url)
}

func TestExtractData(t *testing.T) {
	data := `<html>
		<span>Sector</span>
		<span>TestSector</span>
		<span>Industry</span>
		<span>Test Industry</span>
	</html>
	`
	result := extractData("0700.HK", data, nil)
	assert.Equal(t, "TestSector", result.data["sector"])
	assert.Equal(t, "Test Industry", result.data["industry"])
	assert.Nil(t, result.err)
}

func TestExtractDataWithError(t *testing.T) {
	err := errors.New("Create Error")
	result := extractData("0700.HK", "", err)
	assert.Equal(t, "0700.HK", result.ticker)
	assert.Equal(t, "", result.data["sector"])
	assert.Equal(t, "", result.data["industry"])
	assert.Same(t, err, result.err)
}

func TestExtractDataWithNoSectorOrIndustry(t *testing.T) {
	data := `<html>
		<span></span>
	</html>
	`
	result := extractData("0700.HK", data, nil)
	assert.Equal(t, "", result.data["sector"])
	assert.Equal(t, "", result.data["industry"])
	assert.Nil(t, result.err)
}

func TestExtractDataWithSectorOnly(t *testing.T) {
	data := `<html>
		<span>Sector</span>
		<span>TestSector</span>
	</html>
	`
	result := extractData("0700.HK", data, nil)
	assert.Equal(t, "TestSector", result.data["sector"])
	assert.Equal(t, "", result.data["industry"])
	assert.Nil(t, result.err)
}

func TestExtractDataWithIndustryOnly(t *testing.T) {
	data := `<html>
		<span>Industry</span>
		<span>Test Industry</span>
	</html>
	`
	result := extractData("0700.HK", data, nil)
	assert.Equal(t, "", result.data["sector"])
	assert.Equal(t, "Test Industry", result.data["industry"])
	assert.Nil(t, result.err)
}
