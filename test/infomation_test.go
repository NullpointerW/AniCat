package test

import (
	"fmt"
	"strings"
	"testing"

	I "github.com/NullpointerW/anicat/crawl/information"
	util "github.com/NullpointerW/anicat/utils"
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
	tip, err := I.BgmTVInfoScrape(333979)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	for k, v := range tip {
		fmt.Println(k)
		fmt.Println(v)
	}
}

func TestTMDB(t *testing.T) {
	_, d, e := I.FloderSearch(I.TMDB_TYP_TV, "凉宫春日的忧郁")
	if e != nil {
		t.Error(e)
		t.FailNow()
	}
	pd, _ := util.ParseShort02Time(strings.ReplaceAll(d, " ", ""))
	fmt.Println(pd)
}

func TestBgmiSearchApi(t *testing.T) {
	sid, err := I.BgmiApiSearch("无职转生～到了异世界就拿出真本事～ 第2部分")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(sid)
}
