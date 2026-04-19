package selector

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	"gopkg.in/yaml.v3"
)

type Entry struct {
	XPath    string `yaml:"xpath"`
	Validate string `yaml:"validate"`
}

// Registry: site -> key -> Entry
type Registry map[string]map[string]Entry

// LLMConfig holds provider-agnostic LLM settings for the healer.
type LLMConfig struct {
	Style   string // "anthropic" | "openai" (OpenAI-compatible)
	APIKey  string
	Model   string
	BaseURL string // only for openai-compatible providers
}

var (
	mu        sync.RWMutex
	reg       Registry
	filePath  string
	llmConf   LLMConfig
)

// SetLLMConfig stores the LLM configuration for use by the healer.
func SetLLMConfig(c LLMConfig) {
	mu.Lock()
	llmConf = c
	mu.Unlock()
}

// TriggerHeal runs Heal asynchronously if an API key is configured.
func TriggerHeal(site, key, fetchURL string) {
	mu.RLock()
	c := llmConf
	mu.RUnlock()
	if c.APIKey == "" {
		return
	}
	go Heal(site, key, fetchURL, c)
}

func Load(p string) error {
	b, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("selector: read %s: %w", p, err)
	}
	var r Registry
	if err := yaml.Unmarshal(b, &r); err != nil {
		return fmt.Errorf("selector: parse: %w", err)
	}
	mu.Lock()
	reg = r
	filePath = p
	mu.Unlock()
	return nil
}

// Get returns the XPath for site+key, empty string if not found.
func Get(site, key string) string {
	mu.RLock()
	defer mu.RUnlock()
	if s, ok := reg[site]; ok {
		if e, ok := s[key]; ok {
			return e.XPath
		}
	}
	return ""
}

// GetEntry returns the full Entry including validate pattern.
func GetEntry(site, key string) (Entry, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if s, ok := reg[site]; ok {
		if e, ok := s[key]; ok {
			return e, true
		}
	}
	return Entry{}, false
}

// Validate checks value against the entry's validate regex.
// Returns true if no validate pattern is set.
func Validate(site, key, value string) bool {
	e, ok := GetEntry(site, key)
	if !ok || e.Validate == "" {
		return true
	}
	matched, _ := regexp.MatchString(e.Validate, value)
	return matched
}

// Update writes a new XPath for site+key to memory and persists to disk.
func Update(site, key, xpath string) error {
	mu.Lock()
	defer mu.Unlock()
	if reg == nil {
		reg = make(Registry)
	}
	if reg[site] == nil {
		reg[site] = make(map[string]Entry)
	}
	e := reg[site][key]
	e.XPath = xpath
	reg[site][key] = e
	b, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("selector: marshal: %w", err)
	}
	return os.WriteFile(filePath, b, 0644)
}
