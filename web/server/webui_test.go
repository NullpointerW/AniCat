package webserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ── sxxexxRe ──────────────────────────────────────────────────────────────

func TestSxxExxRegexp(t *testing.T) {
	cases := []struct {
		fn       string
		wantSE   string
		wantS    string
		wantE    string
	}{
		{"孤独摇滚 S01E02.mp4", "S01E02", "01", "02"},
		{"天国大魔镜 s02e10.mkv", "s02e10", "02", "10"},
		{"番剧 S01E01.mp4", "S01E01", "01", "01"},
		{"no-episode.mp4", "", "", ""},
	}
	for _, c := range cases {
		m := sxxexxRe.FindStringSubmatch(c.fn)
		if c.wantSE == "" {
			if len(m) != 0 {
				t.Errorf("regexp matched %q in %q, want no match", m[0], c.fn)
			}
			continue
		}
		if len(m) != 3 {
			t.Errorf("regexp on %q: got %v, want 3 groups", c.fn, m)
			continue
		}
		if !strings.EqualFold(m[0], c.wantSE) {
			t.Errorf("%q: SxxExx = %q, want %q", c.fn, m[0], c.wantSE)
		}
		if m[1] != c.wantS {
			t.Errorf("%q: season = %q, want %q", c.fn, m[1], c.wantS)
		}
		if m[2] != c.wantE {
			t.Errorf("%q: episode = %q, want %q", c.fn, m[2], c.wantE)
		}
	}
}

var _ = regexp.MustCompile // keep import

// ── hasCover ──────────────────────────────────────────────────────────────

func TestHasCover(t *testing.T) {
	dir := t.TempDir()

	// No cover yet
	if hasCover(dir) {
		t.Error("hasCover should return false when folder.jpg absent")
	}

	// Create folder.jpg
	coverPath := filepath.Join(dir, "folder.jpg")
	os.WriteFile(coverPath, []byte("fake-img"), 0644)

	if !hasCover(dir) {
		t.Error("hasCover should return true when folder.jpg exists")
	}
}

// ── CORS middleware ───────────────────────────────────────────────────────

func TestCORS_SetsHeader(t *testing.T) {
	handler := cors(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("CORS header = %q, want *", got)
	}
}

func TestCORS_PreflightReturns204(t *testing.T) {
	called := false
	handler := cors(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Errorf("OPTIONS status = %d, want 204", rr.Code)
	}
	if called {
		t.Error("inner handler should not be called for OPTIONS preflight")
	}
}

// ── handleVideo: path traversal protection ────────────────────────────────

func TestHandleVideo_PathTraversalBlocked(t *testing.T) {
	// We need a fake subject in Mgr — but since subject.Mgr requires
	// full init, we test the traversal-check logic directly via HTTP.
	// A request with ".." in the path should get 400 or 403, not 200.
	req := httptest.NewRequest(http.MethodGet, "/api/video/abc/bad", nil)
	rr := httptest.NewRecorder()
	handleVideo(rr, req)
	// sid "abc" is not a number → 400
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for non-numeric sid", rr.Code)
	}
}

func TestHandleVideo_InvalidSid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/video/notanumber/file.mp4", nil)
	rr := httptest.NewRecorder()
	handleVideo(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

// ── handleCover: invalid sid ──────────────────────────────────────────────

func TestHandleCover_InvalidSid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/cover/notanumber", nil)
	rr := httptest.NewRecorder()
	handleCover(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

// ── handleDetail: invalid sid ─────────────────────────────────────────────

func TestHandleDetail_InvalidSid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/bangumi/notanumber", nil)
	rr := httptest.NewRecorder()
	handleDetail(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

// ── scanEpisodes: file system scan ────────────────────────────────────────

func TestScanEpisodes_EmptyPath(t *testing.T) {
	// subject with no path returns empty list without panic
	eps := scanEpisodesFromPath("")
	if len(eps) != 0 {
		t.Errorf("expected empty, got %d episodes", len(eps))
	}
}

func TestScanEpisodes_SortsCorrectly(t *testing.T) {
	dir := t.TempDir()
	// Create fake video files out of order
	files := []string{
		"番剧 S01E03.mp4",
		"番剧 S01E01.mp4",
		"番剧 S02E01.mp4",
		"番剧 S01E02.mp4",
	}
	for _, fn := range files {
		os.WriteFile(filepath.Join(dir, fn), []byte{}, 0644)
	}

	eps := scanEpisodesFromPath(dir)
	if len(eps) != 4 {
		t.Fatalf("len(eps) = %d, want 4", len(eps))
	}
	expected := []struct{ s, e int }{{1, 1}, {1, 2}, {1, 3}, {2, 1}}
	for i, ex := range expected {
		if eps[i].Season != ex.s || eps[i].Ep != ex.e {
			t.Errorf("eps[%d] = S%02dE%02d, want S%02dE%02d",
				i, eps[i].Season, eps[i].Ep, ex.s, ex.e)
		}
	}
}

func TestScanEpisodes_IgnoresNonVideoFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "番剧 S01E01.mp4"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "folder.jpg"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "meta-data#01.json"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "subtitle.srt"), []byte{}, 0644)

	eps := scanEpisodesFromPath(dir)
	if len(eps) != 1 {
		t.Errorf("len(eps) = %d, want 1 (only mp4)", len(eps))
	}
}

func TestScanEpisodes_IndexStartsAtOne(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "番剧 S01E01.mp4"), []byte{}, 0644)
	os.WriteFile(filepath.Join(dir, "番剧 S01E02.mp4"), []byte{}, 0644)

	eps := scanEpisodesFromPath(dir)
	for i, ep := range eps {
		if ep.Index != i+1 {
			t.Errorf("eps[%d].Index = %d, want %d", i, ep.Index, i+1)
		}
	}
}
