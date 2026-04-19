package subject

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/NullpointerW/anicat/downloader/rss"
	"github.com/NullpointerW/anicat/log"
)

// FilterVerb holds precompiled regexps for RSS item filtering.
// Contain rules use AND logic; Exclusion rules use OR-reject logic.
// Precompilation avoids repeated regexp.Compile on every filter call.
type FilterVerb struct {
	// Runtime fields: precompiled regexps, not persisted.
	contain   []*regexp.Regexp
	exclusion []*regexp.Regexp
	// Persisted fields: raw pattern strings for JSON reload.
	RawContain   []string `json:"contain"`
	RawExclusion []string `json:"exclusion"`
}

// BuildFilterVerb constructs a FilterVerb from raw pattern strings.
func BuildFilterVerb(contain, exclusion []string) *FilterVerb {
	f := &FilterVerb{
		RawContain:   contain,
		RawExclusion: exclusion,
	}
	f.contain = compileRegs(contain)
	f.exclusion = compileRegs(exclusion)
	return f
}

// BuildFilterVerbSingle constructs a FilterVerb from a single regex pair.
func BuildFilterVerbSingle(contain, exclusion string) *FilterVerb {
	return BuildFilterVerb([]string{contain}, []string{exclusion})
}

// Restore recompiles regexps after JSON unmarshal.
func (f *FilterVerb) Restore() {
	f.contain = compileRegs(f.RawContain)
	f.exclusion = compileRegs(f.RawExclusion)
}

// Filter returns an rss.FilterFunc backed by this FilterVerb.
func (f *FilterVerb) Filter() rss.FilterFunc {
	return func(n string) bool {
		return matchRegs(n, f.contain, f.exclusion)
	}
}

// compileRegs compiles pattern strings into []*regexp.Regexp.
// Commas within a pattern become | (OR within a group); all patterns are case-insensitive.
func compileRegs(patterns []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			continue
		}
		p = "(?i)" + strings.ReplaceAll(p, ",", "|")
		re, err := regexp.Compile(p)
		if err != nil {
			log.Error(log.Struct{"err", err, "pattern", p}, "filter: regexp compile failed, skipping")
			continue
		}
		out = append(out, re)
	}
	return out
}

// matchRegs applies precompiled regexps: all contain rules must match (AND),
// any exclusion rule match rejects the string.
func matchRegs(s string, contains, exclusions []*regexp.Regexp) bool {
	for _, re := range contains {
		matched := re.MatchString(s)
		log.Debug(log.Struct{"containRegexp", re.String(), "matchingString", s, "matched", matched})
		if !matched {
			return false
		}
	}
	for _, re := range exclusions {
		if re.MatchString(s) {
			log.Debug(log.Struct{"exclusionRegexp", re.String(), "matchingString", s, "matched", true})
			return false
		}
	}
	return true
}

// BuildFilterPerlReg builds a single perl-style lookahead regex for qBittorrent AutoDL rules.
func BuildFilterPerlReg(vbs []string) string {
	if len(vbs) == 0 {
		return ""
	}
	const tmp = `(?=.*?%s)`
	reg := "(?i)"
	for _, ct := range vbs {
		vb := "(" + strings.ReplaceAll(ct, ",", "|") + ")"
		reg += fmt.Sprintf(tmp, vb)
	}
	return reg
}

// BuildFilterRegs compiles pattern strings into []*regexp.Regexp for repeated use.
func BuildFilterRegs(patterns []string) []*regexp.Regexp {
	return compileRegs(patterns)
}

// FilterWithRegs filters s against precompiled regexp slices.
func FilterWithRegs(s string, contains, exclusions []*regexp.Regexp) bool {
	return matchRegs(s, contains, exclusions)
}

// FilterWithReg checks s against a single contain and exclusion pattern string.
func FilterWithReg(s, contain, exclusion string) bool {
	var c, e []*regexp.Regexp
	if contain != "" {
		c = compileRegs([]string{contain})
	}
	if exclusion != "" {
		e = compileRegs([]string{exclusion})
	}
	return matchRegs(s, c, e)
}

// FilterWithCustomReg applies per-subject regex filters from an Extra.
func FilterWithCustomReg(s string, ex Extra) bool {
	return FilterWithReg(s, ex.RssOption.MustContain, ex.RssOption.MustNotContain)
}

// FilterWithCustom applies per-subject keyword filters from an Extra (splits on whitespace).
func FilterWithCustom(s string, ex Extra) bool {
	return FilterWithRegs(s,
		compileRegs(strings.Fields(ex.RssOption.MustContain)),
		compileRegs(strings.Fields(ex.RssOption.MustNotContain)),
	)
}
