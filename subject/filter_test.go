package subject

import (
	"encoding/json"
	"strings"
	"testing"
)

// ── BuildFilterVerb / Filter ──────────────────────────────────────────────

func TestBuildFilterVerb_ContainMatch(t *testing.T) {
	f := BuildFilterVerb([]string{"1080p"}, nil)
	fn := f.Filter()
	cases := []struct {
		input string
		want  bool
	}{
		{"[Lilith-Raws] 孤独摇滚 S01E01 [1080p].mkv", true},
		{"[Lilith-Raws] 孤独摇滚 S01E01 [720p].mkv", false},
		{"[Lilith-Raws] 孤独摇滚 S01E01 [1080P].mkv", true}, // case-insensitive
	}
	for _, c := range cases {
		if got := fn(c.input); got != c.want {
			t.Errorf("Filter(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestBuildFilterVerb_ExclusionReject(t *testing.T) {
	f := BuildFilterVerb(nil, []string{"720p", "480p"})
	fn := f.Filter()
	cases := []struct {
		input string
		want  bool
	}{
		{"[Sub] 番剧 S01E01 [1080p].mkv", true},
		{"[Sub] 番剧 S01E01 [720p].mkv", false},
		{"[Sub] 番剧 S01E01 [480P].mkv", false},
	}
	for _, c := range cases {
		if got := fn(c.input); got != c.want {
			t.Errorf("Filter(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestBuildFilterVerb_ContainANDLogic(t *testing.T) {
	f := BuildFilterVerb([]string{"1080p", "简体"}, nil)
	fn := f.Filter()
	cases := []struct {
		input string
		want  bool
	}{
		{"[Sub] 番剧 [1080p][简体].mkv", true},
		{"[Sub] 番剧 [1080p][繁体].mkv", false},
		{"[Sub] 番剧 [720p][简体].mkv", false},
	}
	for _, c := range cases {
		if got := fn(c.input); got != c.want {
			t.Errorf("Filter(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestBuildFilterVerb_CommaAsOR(t *testing.T) {
	f := BuildFilterVerb([]string{"简体,繁体"}, nil)
	fn := f.Filter()
	cases := []struct {
		input string
		want  bool
	}{
		{"番剧 [1080p][简体].mkv", true},
		{"番剧 [1080p][繁体].mkv", true},
		{"番剧 [1080p][日语].mkv", false},
	}
	for _, c := range cases {
		if got := fn(c.input); got != c.want {
			t.Errorf("Filter(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestBuildFilterVerb_NilFilter(t *testing.T) {
	f := BuildFilterVerb(nil, nil)
	fn := f.Filter()
	if !fn("anything goes") {
		t.Error("empty FilterVerb should pass all inputs")
	}
}

func TestBuildFilterVerb_EmptyStringSkipped(t *testing.T) {
	f := BuildFilterVerb([]string{"", "1080p", ""}, nil)
	fn := f.Filter()
	if !fn("[Test] 番剧 [1080p].mkv") {
		t.Error("empty pattern strings should be skipped, valid rule should still apply")
	}
}

// ── JSON round-trip ───────────────────────────────────────────────────────

func TestFilterVerb_RestoreAfterJSON(t *testing.T) {
	original := BuildFilterVerb([]string{"1080p"}, []string{"720p"})
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	restored := &FilterVerb{}
	if err := json.Unmarshal(data, restored); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	restored.Restore()
	fn := restored.Filter()
	if !fn("[Test] S01E01 [1080p].mkv") {
		t.Error("restored filter should match 1080p")
	}
	if fn("[Test] S01E01 [720p].mkv") {
		t.Error("restored filter should reject 720p")
	}
}

func TestFilterVerb_JSONPreservesRawFields(t *testing.T) {
	f := BuildFilterVerb([]string{"1080p", "简体"}, []string{"720p"})
	data, _ := json.Marshal(f)
	f2 := &FilterVerb{}
	json.Unmarshal(data, f2)
	if len(f2.RawContain) != 2 {
		t.Errorf("RawContain len = %d, want 2", len(f2.RawContain))
	}
	if len(f2.RawExclusion) != 1 {
		t.Errorf("RawExclusion len = %d, want 1", len(f2.RawExclusion))
	}
}

// ── BuildFilterVerbSingle ─────────────────────────────────────────────────

func TestBuildFilterVerbSingle(t *testing.T) {
	f := BuildFilterVerbSingle("1080p", "720p")
	fn := f.Filter()
	if !fn("番剧 [1080p]") {
		t.Error("should match 1080p")
	}
	if fn("番剧 [720p]") {
		t.Error("should reject 720p")
	}
}

// ── FilterWithRegs ────────────────────────────────────────────────────────

func TestFilterWithRegs(t *testing.T) {
	cases := []struct {
		s         string
		contains  []string
		exclusion []string
		want      bool
	}{
		{"番剧 [1080p][简体]", []string{"1080p"}, nil, true},
		{"番剧 [720p][简体]", []string{"1080p"}, nil, false},
		{"番剧 [1080p][繁体]", nil, []string{"繁体"}, false},
		{"番剧 [1080p][简体]", nil, []string{"繁体"}, true},
		{"番剧 [1080p][简体]", []string{"1080p"}, []string{"繁体"}, true},
		{"番剧 [1080p][繁体]", []string{"1080p"}, []string{"繁体"}, false},
		{"anything", nil, nil, true},
	}
	for _, c := range cases {
		got := FilterWithRegs(c.s, BuildFilterRegs(c.contains), BuildFilterRegs(c.exclusion))
		if got != c.want {
			t.Errorf("FilterWithRegs(%q, %v, %v) = %v, want %v",
				c.s, c.contains, c.exclusion, got, c.want)
		}
	}
}

// ── FilterWithReg ─────────────────────────────────────────────────────────

func TestFilterWithReg(t *testing.T) {
	cases := []struct {
		s         string
		contain   string
		exclusion string
		want      bool
	}{
		{"番剧 [1080p]", "1080p", "", true},
		{"番剧 [720p]", "1080p", "", false},
		{"番剧 [1080p][繁体]", "", "繁体", false},
		{"番剧 [1080p][简体]", "", "繁体", true},
		{"番剧 [1080p][简体]", "1080p", "繁体", true},
		{"anything", "", "", true},
	}
	for _, c := range cases {
		got := FilterWithReg(c.s, c.contain, c.exclusion)
		if got != c.want {
			t.Errorf("FilterWithReg(%q, %q, %q) = %v, want %v",
				c.s, c.contain, c.exclusion, got, c.want)
		}
	}
}

// ── BuildFilterPerlReg ────────────────────────────────────────────────────

func TestBuildFilterPerlReg_Empty(t *testing.T) {
	if got := BuildFilterPerlReg(nil); got != "" {
		t.Errorf("expected empty string for nil, got %q", got)
	}
	if got := BuildFilterPerlReg([]string{}); got != "" {
		t.Errorf("expected empty string for empty slice, got %q", got)
	}
}

func TestBuildFilterPerlReg_ProducesNonEmpty(t *testing.T) {
	// BuildFilterPerlReg generates PCRE lookahead patterns for qBittorrent AutoDL rules.
	// Go's regexp (RE2) does NOT support lookaheads — only check structural properties.
	cases := []struct {
		input []string
		check string
	}{
		{[]string{"1080p"}, "1080p"},
		{[]string{"1080p", "简体"}, "简体"},
		{[]string{"简体,繁体"}, "简体|繁体"},
	}
	for _, c := range cases {
		got := BuildFilterPerlReg(c.input)
		if got == "" {
			t.Errorf("BuildFilterPerlReg(%v) returned empty string", c.input)
		}
		if !strings.Contains(got, c.check) {
			t.Errorf("BuildFilterPerlReg(%v) = %q, want to contain %q", c.input, got, c.check)
		}
		if !strings.Contains(got, "(?=") {
			t.Errorf("BuildFilterPerlReg(%v) = %q, missing lookahead (?=", c.input, got)
		}
		if !strings.HasPrefix(got, "(?i)") {
			t.Errorf("BuildFilterPerlReg(%v) = %q, missing (?i) prefix", c.input, got)
		}
	}
}

func TestBuildFilterPerlReg_TwoWordsProduceTwoGroups(t *testing.T) {
	got := BuildFilterPerlReg([]string{"1080p", "简体"})
	count := strings.Count(got, "(?=")
	if count != 2 {
		t.Errorf("expected 2 lookahead groups, got %d in %q", count, got)
	}
}
