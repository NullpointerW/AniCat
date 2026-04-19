package llmparser

// EpisodeParser extracts a 2-digit episode number string (e.g. "04") from a raw filename.
type EpisodeParser interface {
	Parse(filename string) (episode string, err error)
}

// Fallback is called by rename.CaptureEpisNum when all regex patterns fail.
// Nil means LLM parsing is disabled.
var Fallback EpisodeParser
