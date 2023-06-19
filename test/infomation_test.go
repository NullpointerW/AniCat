package test

import (
	"fmt"
	"testing"

	I "github.com/NullpointerW/anicat/crawl/information"
)

func TestInfoSearch(t *testing.T) {
	I.InfoPageScrape("凉宫春日的消失")
}

func TestInfoScraping(t *testing.T) {
	tip, err := I.Scrape("铃芽之旅")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for k, v := range tip {
		fmt.Println(k)
		fmt.Println(v)
	}
}
func TestBgmTVInfoScrape(t *testing.T) {
	tip, err := I.BgmTVInfoScrape(485)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for k, v := range tip {
		fmt.Println(k)
		fmt.Println(v)
	}
}
