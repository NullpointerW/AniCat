package test

import (
	"fmt"
	I "github.com/NullpointerW/mikanani/crawl/information"
	"testing"
)

func TestInfoSearch(t *testing.T) {
	I.InfoPageScrape("凉宫春日的消失")
}

func TestInfoScraping(t *testing.T) {
	tip, err := I.InfoScrape("铃芽之旅")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for k, v := range tip {
		fmt.Println(k)
		fmt.Println(v)
	}
}
