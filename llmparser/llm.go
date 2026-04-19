package llmparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/log"
)

const (
	styleAnthropic = "anthropic"
	styleOpenAI    = "openai"

	defaultAnthropicBase = "https://api.anthropic.com"
	defaultOpenAIBase    = "https://api.openai.com"
)

type llmParser struct {
	apiKey  string
	model   string
	baseURL string
	style   string // styleAnthropic or styleOpenAI
	client  *http.Client
}

const episodePrompt = "Extract the episode number from this anime filename. " +
	"Return ONLY a JSON object: {\"episode\": \"04\"}. " +
	"Use a 2-digit zero-padded string. " +
	"If you cannot determine the episode number, return {\"episode\": \"\"}.\n\nFilename: %s"

func (p *llmParser) Parse(filename string) (string, error) {
	prompt := fmt.Sprintf(episodePrompt, filename)
	body, _ := json.Marshal(map[string]any{
		"model":      p.model,
		"max_tokens": 64,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	})

	var url string
	if p.style == styleAnthropic {
		url = p.baseURL + "/v1/messages"
	} else {
		url = p.baseURL + "/v1/chat/completions"
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llmparser: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if p.style == styleAnthropic {
		req.Header.Set("x-api-key", p.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	} else {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("llmparser: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llmparser: unexpected status %d", resp.StatusCode)
	}

	text, err := extractText(p.style, resp)
	if err != nil {
		return "", err
	}

	var result struct {
		Episode string `json:"episode"`
	}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return "", fmt.Errorf("llmparser: parse json %q: %w", text, err)
	}
	if result.Episode == "" {
		return "", fmt.Errorf("llmparser: could not identify episode in %q", filename)
	}

	log.Info(log.Struct{"file", filename, "episode", result.Episode}, "llmparser: episode identified")
	return result.Episode, nil
}

func extractText(style string, resp *http.Response) (string, error) {
	if style == styleAnthropic {
		var r struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return "", fmt.Errorf("llmparser: decode response: %w", err)
		}
		if len(r.Content) == 0 {
			return "", fmt.Errorf("llmparser: empty response")
		}
		return strings.TrimSpace(r.Content[0].Text), nil
	}

	var r struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("llmparser: decode response: %w", err)
	}
	if len(r.Choices) == 0 {
		return "", fmt.Errorf("llmparser: empty response")
	}
	return strings.TrimSpace(r.Choices[0].Message.Content), nil
}

func init() {
	if CFG.SrvCTL {
		return
	}
	cfg := CFG.Env.LLMParser
	if cfg.APIKey == "" {
		return
	}

	style := cfg.Style
	if style == "" {
		style = styleOpenAI
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		if style == styleAnthropic {
			baseURL = defaultAnthropicBase
		} else {
			baseURL = defaultOpenAIBase
		}
	}

	Fallback = &llmParser{
		apiKey:  cfg.APIKey,
		model:   cfg.Model,
		baseURL: strings.TrimRight(baseURL, "/"),
		style:   style,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	log.Info(log.Struct{"style", style, "model", cfg.Model, "base_url", baseURL}, "llmparser: enabled")
}
