package test

import (
	"fmt"
	"testing"

	R "github.com/NullpointerW/mikanani/crawl/resource"
)

func TestCrwal(t *testing.T) {
	n := "凉宫春日"
	url, _, isrss, err := R.Scrape(n)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("name:%s|is rss_resource :%v|url:%s\n", n, isrss, url)
	fmt.Println("====================================")
	n = "lycoris Recoil"
	url, bgm, isrss, err := R.Scrape(n)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Printf("name:%s|is rss_resource :%v|url:%s\n", n, isrss, url)
	fmt.Println("bgm_url:" + bgm)
}
