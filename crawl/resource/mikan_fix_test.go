package resource

import (
	"strings"
	"testing"
)

// Fix: Mikan added a checkbox column; td indices shifted +1.
// Tests the flat search table path (searchls=true → uses fnTemp/szTemp/uptTemp).
func TestMikanFlatListColumns(t *testing.T) {
	res, err := ListScrape("葬送的芙莉莲", Ls, true)
	if err != nil {
		t.Errorf("ListScrape (flat) failed: %v", err)
		t.FailNow()
	}
	items, ok := res.([]Item)
	if !ok || len(items) == 0 {
		t.Error("expected non-empty item list")
		t.FailNow()
	}
	first := items[0]
	t.Logf("name=%q  size=%q  uptime=%q", first.Name, first.Size, first.UpdateTime)

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

// Tests the RSS group path (searchls=false → scrapeRssList uses itnExp/szExp/uptExp).
func TestMikanRssGroupColumns(t *testing.T) {
	res, err := ListScrape("葬送的芙莉莲", Ls, false)
	if err != nil {
		t.Errorf("ListScrape (rss group) failed: %v", err)
		t.FailNow()
	}
	rgs, ok := res.([]RssGroup)
	if !ok || len(rgs) == 0 {
		t.Error("expected non-empty RssGroup list")
		t.FailNow()
	}
	checked := false
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
		checked = true
		break
	}
	if !checked {
		t.Error("all RssGroups have empty Items — trsTemp XPath not matching (episode-table wrapper?)")
	}
}

// Fix: mglinkTemp in Scrape changed td[1]/a[2] → td[2]/a[2]
func TestMikanMagnetLink(t *testing.T) {
	url, _, _, err := Scrape("葬送的芙莉莲", Option{Index: 1})
	if err != nil {
		t.Errorf("Scrape (magnet) failed: %v", err)
		t.FailNow()
	}
	if !strings.HasPrefix(url, "magnet:?") {
		t.Errorf("expected magnet URL, got: %q — td[2]/a[2] xpath may be wrong", url)
	}
	if len(url) > 80 {
		t.Logf("magnet=%s...", url[:80])
	} else {
		t.Logf("magnet=%s", url)
	}
}
