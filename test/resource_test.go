package test

import (
	"fmt"
	"strings"
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

// Fix: Mikan added a checkbox column; td[1]=name → td[2]=name, td[2]=size → td[3], td[3]=uptime → td[4]
// Tests the flat search table path (searchls=true → uses fnTemp/szTemp/uptTemp in ListScrape)
func TestMikanFlatListColumns(t *testing.T) {
	res, err := R.ListScrape("葬送的芙莉莲", R.Ls, true)
	if err != nil {
		t.Errorf("ListScrape (flat) failed: %v", err)
		t.FailNow()
	}
	items, ok := res.([]R.Item)
	if !ok || len(items) == 0 {
		t.Error("expected non-empty item list")
		t.FailNow()
	}
	first := items[0]
	t.Logf("name=%q  size=%q  uptime=%q", first.Name, first.Size, first.UpdateTime)

	// If td indices are wrong, Name gets the size string and Size gets the uptime
	if strings.Contains(first.Name, "MB") || strings.Contains(first.Name, "GB") {
		t.Errorf("Name looks like a size string — td index still off: %q", first.Name)
	}
	if first.Name == "" {
		t.Error("Name is empty — td index or xpath mismatch")
	}
	if !strings.Contains(first.Size, "MB") && !strings.Contains(first.Size, "GB") {
		t.Errorf("Size field doesn't look like a size: %q (expected MB/GB)", first.Size)
	}
}

// Tests the RSS group path (searchls=false → scrapeRssList uses itnExp/szExp/uptExp)
func TestMikanRssGroupColumns(t *testing.T) {
	res, err := R.ListScrape("葬送的芙莉莲", R.Ls, false)
	if err != nil {
		t.Errorf("ListScrape (rss group) failed: %v", err)
		t.FailNow()
	}
	rgs, ok := res.([]R.RssGroup)
	if !ok || len(rgs) == 0 {
		t.Error("expected non-empty RssGroup list")
		t.FailNow()
	}
	for _, rg := range rgs {
		if len(rg.Items) == 0 {
			continue
		}
		first := rg.Items[0]
		t.Logf("[%s] name=%q  size=%q  uptime=%q", rg.Name, first.Name, first.Size, first.UpdateTime)
		if strings.Contains(first.Name, "MB") || strings.Contains(first.Name, "GB") {
			t.Errorf("group %q: Name looks like a size string — td index still off: %q", rg.Name, first.Name)
		}
		if first.Name == "" {
			t.Errorf("group %q: Name is empty", rg.Name)
		}
		if !strings.Contains(first.Size, "MB") && !strings.Contains(first.Size, "GB") {
			t.Errorf("group %q: Size doesn't look like a size: %q", rg.Name, first.Size)
		}
		break
	}
}

// Fix: mglinkTemp in Scrape changed td[1]/a[2] → td[2]/a[2]
func TestMikanMagnetLink(t *testing.T) {
	url, _, _, err := R.Scrape("葬送的芙莉莲", R.Option{Index: 1})
	if err != nil {
		t.Errorf("Scrape (magnet) failed: %v", err)
		t.FailNow()
	}
	if !strings.HasPrefix(url, "magnet:?") {
		t.Errorf("expected magnet URL, got: %q — td[2]/a[2] xpath may be wrong", url)
	}
	t.Logf("magnet=%s...", url[:min(len(url), 80)])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
