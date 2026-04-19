package util

import "testing"

// ── IsVideofile ───────────────────────────────────────────────────────────

func TestIsVideofile(t *testing.T) {
	cases := []struct {
		fn   string
		want bool
	}{
		{"孤独摇滚 S01E01.mp4", true},
		{"孤独摇滚 S01E01.mkv", true},
		{"孤独摇滚 S01E01.MKV", true},   // uppercase extension
		{"孤独摇滚 S01E01.MP4", true},
		{"孤独摇滚 S01E01.avi", true},
		{"孤独摇滚 S01E01.mov", true},
		{"cover.jpg", false},
		{"meta-data#01.json", false},
		{"subtitle.srt", false},
		{"subtitle.ass", false},
		{"", false},
		{"noextension", false},
		{"孤独摇滚 S01E01.mp4.bak", false},
	}
	for _, c := range cases {
		if got := IsVideofile(c.fn); got != c.want {
			t.Errorf("IsVideofile(%q) = %v, want %v", c.fn, got, c.want)
		}
	}
}

// ── IsSubtitleFile ────────────────────────────────────────────────────────

func TestIsSubtitleFile(t *testing.T) {
	cases := []struct {
		fn   string
		want bool
	}{
		{"subtitle.srt", true},
		{"subtitle.SRT", true},
		{"subtitle.ass", true},
		{"subtitle.sub", true},
		{"video.mp4", false},
		{"cover.jpg", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsSubtitleFile(c.fn); got != c.want {
			t.Errorf("IsSubtitleFile(%q) = %v, want %v", c.fn, got, c.want)
		}
	}
}

// ── IsJsonFile ────────────────────────────────────────────────────────────

func TestIsJsonFile(t *testing.T) {
	cases := []struct {
		fn   string
		want bool
	}{
		{"meta-data#01.json", true},
		{"config.JSON", true},
		{"video.mp4", false},
		{"", false},
	}
	for _, c := range cases {
		if got := IsJsonFile(c.fn); got != c.want {
			t.Errorf("IsJsonFile(%q) = %v, want %v", c.fn, got, c.want)
		}
	}
}

// ── FileSeparatorConv ─────────────────────────────────────────────────────

func TestFileSeparatorConv(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{`C:\Users\test\bangumi`, "C:/Users/test/bangumi"},
		{"/home/user/bangumi", "/home/user/bangumi"},
		{`path\to\file`, "path/to/file"},
		{"", ""},
		{"no-backslash", "no-backslash"},
	}
	for _, c := range cases {
		if got := FileSeparatorConv(c.input); got != c.want {
			t.Errorf("FileSeparatorConv(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// ── TrimExtensionAndGetEpi ────────────────────────────────────────────────

func TestTrimExtensionAndGetEpi(t *testing.T) {
	cases := []struct {
		fn   string
		want string
	}{
		{"孤独摇滚！S01E02.mp4", "S01E02"},
		{"天国大魔镜 S01E10.mkv", "S01E10"},
		{"番剧 S02E01.mp4", "S02E01"},
		// last 6 chars after trim
		{"ABCDEFGHIJ S01E99.mp4", "S01E99"},
		// short string — return full trimmed
		{"ab.mp4", "ab"},
		{"a.mp4", "a"},
	}
	for _, c := range cases {
		if got := TrimExtensionAndGetEpi(c.fn); got != c.want {
			t.Errorf("TrimExtensionAndGetEpi(%q) = %q, want %q", c.fn, got, c.want)
		}
	}
}

// ── Bug regression: strings shorter than 6 chars after trim must not panic ──

func TestTrimExtensionAndGetEpi_ShortNoPanic(t *testing.T) {
	shorts := []string{"a.mp4", "ab.mkv", "abc.avi", "abcd.mov", "abcde.mp4"}
	for _, fn := range shorts {
		// Must not panic
		got := TrimExtensionAndGetEpi(fn)
		if got == "" {
			t.Errorf("TrimExtensionAndGetEpi(%q) returned empty string unexpectedly", fn)
		}
	}
}

// ── IsRegexp ──────────────────────────────────────────────────────────────

func TestIsRegexp(t *testing.T) {
	cases := []struct {
		pattern string
		want    bool
	}{
		{"1080p", true},
		{"[valid]", true},
		{"(?i)hello", true},
		{"[invalid", false},
		{"(unclosed", false},
		{"", true}, // empty string is valid regexp
	}
	for _, c := range cases {
		if got := IsRegexp(c.pattern); got != c.want {
			t.Errorf("IsRegexp(%q) = %v, want %v", c.pattern, got, c.want)
		}
	}
}
