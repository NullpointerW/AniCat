package test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/NullpointerW/anicat/downloader/rss"
	"github.com/NullpointerW/anicat/subject"
)

func TestRss(t *testing.T) {
	s := subject.Subject{}
	s.SubjId = 379639
	as, err := rss.GetMatchedArts(s.RssPath())
	if err != nil {
		t.Error(as)
		t.FailNow()
	}
	fmt.Println("macthed lens", len(as))
}

func TestRssItem(t *testing.T) {
	it, _ := rss.GetItems("anicat@subj-115908")
	fmt.Println(*it)
}

func TestFeed(t *testing.T) {
	p := rss.Parser{Feed: "https://mikanani.me/RSS/Bangumi?bangumiId=2549&subgroupid=534"}
	_, bgmiUrl, err := p.GetTitleAndLink()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bgmiUrl)
}

func TestBgmTVReg(t *testing.T) {
	str := "http://mikanani.me/RSS/Bangumi?bangumiId=2549&subgroupid=534"
	re := regexp.MustCompile(`bangumiId=(\d+)`)
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		season := match[1]
		fmt.Printf("匹配到的id：%s\n", season)
	} else {
		fmt.Println("未匹配到part")
	}
}
