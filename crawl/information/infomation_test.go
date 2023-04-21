package infomation

import (
	"fmt"
	"testing"
)

func TestInfoSearch(t *testing.T) {
	InfoPageScrape("凉宫春日的消失")
}

func TestInfoScraping(t *testing.T) {
	tip, err := InfoScrape("凉宫春日的忧郁2009")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for k, v := range tip {
		fmt.Println(k)
		fmt.Println(v)
	}
}
