package test

import (
	"fmt"
	"testing"

	"github.com/NullpointerW/anicat/download/rss"
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
