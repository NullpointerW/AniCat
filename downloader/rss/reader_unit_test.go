package rss

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockRSSFeed is a minimal valid RSS feed with 3 items.
const mockRSSFeed = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Test Feed</title>
    <item>
      <title>[Sub] 番剧 S01E01 [1080p].mkv</title>
      <guid>guid-001</guid>
      <enclosure url="http://example.com/s01e01.torrent" type="application/x-bittorrent"/>
    </item>
    <item>
      <title>[Sub] 番剧 S01E02 [1080p].mkv</title>
      <guid>guid-002</guid>
      <enclosure url="http://example.com/s01e02.torrent" type="application/x-bittorrent"/>
    </item>
    <item>
      <title>[Sub] 番剧 S01E03 [720p].mkv</title>
      <guid>guid-003</guid>
      <enclosure url="http://example.com/s01e03.torrent" type="application/x-bittorrent"/>
    </item>
  </channel>
</rss>`

func newMockServer(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(body))
	}))
}

// ── Read: basic ───────────────────────────────────────────────────────────

func TestRead_ReturnsAllItems(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	items, ok, err := r.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !ok {
		t.Fatal("Read() ok = false, want true")
	}
	if len(items) != 3 {
		t.Errorf("len(items) = %d, want 3", len(items))
	}
}

func TestRead_DeduplicatesOnSecondCall(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	r.Read() // first call: consume all 3

	items, ok, err := r.Read()
	if err != nil {
		t.Fatalf("second Read() error: %v", err)
	}
	if ok {
		t.Error("second Read() ok = true, want false (all seen)")
	}
	if len(items) != 0 {
		t.Errorf("second Read() len = %d, want 0", len(items))
	}
}

// ── Read: filter ──────────────────────────────────────────────────────────

func TestRead_FilterPassesOnlyMatching(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	// Only 1080p items should pass
	r2 := NewReader(srv.URL, nil, func(title string) bool {
		for i := range title {
			if i+5 <= len(title) && title[i:i+5] == "1080p" {
				return true
			}
		}
		return false
	})

	items, ok, err := r2.Read()
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if !ok {
		t.Fatal("Read() ok = false, want true")
	}
	if len(items) != 2 {
		t.Errorf("len(items) = %d, want 2 (only 1080p)", len(items))
	}
}

func TestRead_FilteredOutItemsNotMarkedSeen(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	// First read: filter allows only 1080p (guids 001, 002 marked seen; 003 not)
	onlyHD := func(title string) bool {
		for i := range title {
			if i+5 <= len(title) && title[i:i+5] == "1080p" {
				return true
			}
		}
		return false
	}
	r := NewReader(srv.URL, nil, onlyHD)
	r.Read()

	// Second read: remove filter → guid-003 (720p) should now appear
	r2 := NewReader(srv.URL, r.Guids(), nil)
	items, ok, err := r2.Read()
	if err != nil {
		t.Fatalf("second Read() error: %v", err)
	}
	if !ok {
		t.Error("second Read() ok = false, want true (guid-003 was not seen)")
	}
	if len(items) != 1 {
		t.Errorf("len(items) = %d, want 1 (only guid-003 unseen)", len(items))
	}
	if len(items) == 1 && items[0].Guid != "guid-003" {
		t.Errorf("expected guid-003, got %q", items[0].Guid)
	}
}

// ── ReadOne ───────────────────────────────────────────────────────────────

func TestReadOne_ReturnsFirstItem(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	item, ok, err := r.ReadOne()
	if err != nil {
		t.Fatalf("ReadOne() error: %v", err)
	}
	if !ok {
		t.Fatal("ReadOne() ok = false, want true")
	}
	if item.Guid != "guid-001" {
		t.Errorf("Guid = %q, want guid-001", item.Guid)
	}
}

func TestReadOne_ProgressesThroughItems(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	guids := []string{"guid-001", "guid-002", "guid-003"}
	for _, want := range guids {
		item, ok, err := r.ReadOne()
		if err != nil {
			t.Fatalf("ReadOne() error: %v", err)
		}
		if !ok {
			t.Fatalf("ReadOne() ok = false, want true for %q", want)
		}
		if item.Guid != want {
			t.Errorf("Guid = %q, want %q", item.Guid, want)
		}
	}
	// 4th call: all exhausted
	_, ok, err := r.ReadOne()
	if err != nil {
		t.Fatalf("4th ReadOne() error: %v", err)
	}
	if ok {
		t.Error("4th ReadOne() ok = true, want false (all exhausted)")
	}
}

// ── Seek ──────────────────────────────────────────────────────────────────

func TestSeek_DoesNotMarkSeen(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	items1, _, _ := r.Seek()
	items2, _, _ := r.Seek()
	if len(items1) != len(items2) {
		t.Errorf("Seek() not idempotent: first=%d second=%d", len(items1), len(items2))
	}
}

// ── Guids / Undo ──────────────────────────────────────────────────────────

func TestGuids_Snapshot(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	r.Read()

	snap := r.Guids()
	if len(snap) != 3 {
		t.Errorf("Guids() len = %d, want 3", len(snap))
	}
	// Mutating snapshot must not affect internal state
	delete(snap, "guid-001")
	if len(r.Guids()) != 3 {
		t.Error("Guids() snapshot mutation affected internal state")
	}
}

func TestUndo_AllowsReprocessing(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	r.Read() // consume all

	r.Undo("guid-001")
	items, ok, err := r.Read()
	if err != nil {
		t.Fatalf("Read() after Undo error: %v", err)
	}
	if !ok {
		t.Fatal("Read() after Undo: ok = false, want true")
	}
	if len(items) != 1 || items[0].Guid != "guid-001" {
		t.Errorf("after Undo: got %v, want [guid-001]", items)
	}
}

// ── Item fields ───────────────────────────────────────────────────────────

func TestRead_ItemFieldsPopulated(t *testing.T) {
	srv := newMockServer(mockRSSFeed)
	defer srv.Close()

	r := NewReader(srv.URL, nil, nil)
	items, _, _ := r.Read()
	if len(items) == 0 {
		t.Fatal("no items returned")
	}
	it := items[0]
	if it.Guid == "" {
		t.Error("Guid is empty")
	}
	if it.Title == "" {
		t.Error("Title is empty")
	}
	if it.TorrUrl == "" {
		t.Error("TorrUrl is empty")
	}
}
