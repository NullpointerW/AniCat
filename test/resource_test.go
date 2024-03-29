package test

import (
	"fmt"
	"testing"

	R "github.com/NullpointerW/anicat/crawl/resource"
)

func TestCrawl(t *testing.T) {
	n := "小林家的龙女仆"
	url, _, isrss, err := R.Scrape(n, R.Option{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("name:%s|is rss_resource :%v|url:%s\n", n, isrss, url)
	// fmt.Println("====================================")
	// n = "总之就是非常可爱 第二季"
	// url, bgm, isrss, err := R.Scrape(n, R.Option{})
	// if err != nil {
	// 	t.Error(err)
	// 	t.FailNow()
	// }
	// fmt.Printf("name:%s|is rss_resource :%v|url:%s\n", n, isrss, url)
	// fmt.Println("bgm_url:" + bgm)
}

func TestList(t *testing.T) {
	res, err := R.ListScrape("总之就是非常可爱 第二季", R.Ls, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%#+v", res)
}

func TestLsGroup(t *testing.T) {
	res, err := R.ListScrape("总之就是非常可爱 第二季", R.LSGroup, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("%#+v", res)
}

func TestRssOptCrawl(t *testing.T) {
	n := "总之就是非常可爱 第二季"
	url, bgm, isrss, err := R.Scrape(n, R.Option{
		Group: "ANi",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("name:%s|is rss_resource :%v|url:%s\n", n, isrss, url)
	fmt.Println("bgm_url:" + bgm)
}
