package information

import (
	"strings"
	"testing"
)

// Fix: BgmiRoot changed from http:// to https:// (HTTP returned 403)
func TestBgmiApiHttps(t *testing.T) {
	sid, err := BgmiApiSearch("葬送的芙莉莲")
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
	tips, err := BgmTVInfoScrape(400602) // 葬送的芙莉莲
	if err != nil {
		t.Errorf("BgmTVInfoScrape failed — XPath fix may have reverted: %v", err)
		t.FailNow()
	}
	required := []string{SubjName, SubjEpisode, SubjStartTime}
	for _, f := range required {
		v, ok := tips[f]
		if !ok || strings.TrimSpace(v) == "" {
			t.Errorf("missing or empty field %q — infoXpathExp not matching", f)
		} else {
			t.Logf("%-12s = %s", f, v)
		}
	}
	if origin := tips[SubjOriginName]; strings.TrimSpace(origin) == "" {
		t.Errorf("missing origin name — OriginNameXpath not matching")
	} else {
		t.Logf("%-12s = %s", SubjOriginName, origin)
	}
}

// Fix: BgmiApiSearch was comparing URL-encoded searchstr against raw NameCN
func TestBgmiApiNameMatch(t *testing.T) {
	sid, err := BgmiApiSearch("葬送的芙莉莲")
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
