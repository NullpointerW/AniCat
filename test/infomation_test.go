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

// Fix: BgmiRoot changed from http:// to https:// (HTTP returned 403)
func TestBgmiApiHttps(t *testing.T) {
	sid, err := I.BgmiApiSearch("葬送的芙莉莲")
	if err != nil {
		t.Errorf("BGM.TV API unreachable — HTTPS fix may have reverted: %v", err)
		t.FailNow()
	}
	if sid <= 0 {
		t.Errorf("expected positive sid, got %d", sid)
	}
	t.Logf("sid=%d", sid)
}

// Fix: infoXpathExp changed from full path (matched empty div) to //ul[@id='infobox']/li
func TestBgmTVInfoFields(t *testing.T) {
	tips, err := I.BgmTVInfoScrape(400602) // 葬送的芙莉莲
	if err != nil {
		t.Errorf("BgmTVInfoScrape failed — XPath fix may have reverted: %v", err)
		t.FailNow()
	}
	required := []string{I.SubjName, I.SubjEpisode, I.SubjStartTime}
	for _, f := range required {
		v, ok := tips[f]
		if !ok || strings.TrimSpace(v) == "" {
			t.Errorf("missing or empty field %q — infoXpathExp not matching", f)
		} else {
			t.Logf("%-12s = %s", f, v)
		}
	}
	originName := tips[I.SubjOriginName]
	if strings.TrimSpace(originName) == "" {
		t.Errorf("missing origin name — OriginNameXpath not matching")
	} else {
		t.Logf("%-12s = %s", I.SubjOriginName, originName)
	}
}

// Fix: BgmiApiSearch was comparing URL-encoded searchstr against raw NameCN
func TestBgmiApiNameMatch(t *testing.T) {
	sid, err := I.BgmiApiSearch("葬送的芙莉莲")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	const wantSid = 400602
	if sid != wantSid {
		t.Errorf("name matching returned wrong subject: got sid=%d, want %d", sid, wantSid)
	}
	t.Logf("matched sid=%d", sid)
}
