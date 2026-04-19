package selector

import (
	"context"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/NullpointerW/anicat/log"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/antchfx/htmlquery"
)

const (
	maxRounds    = 10
	healerModel  = anthropic.ModelClaudeHaiku4_5
	snippetLimit = 3000
)

// Heal uses an LLM agent loop to repair a broken XPath selector.
// Returns true if a new selector was successfully committed.
func Heal(site, key, fetchURL string, conf LLMConfig) bool {
	log.Info(log.Struct{"site", site, "key", key, "style", conf.Style}, "SelectorHealer: starting repair")
	if conf.Style == "openai" {
		return healWithOpenAI(site, key, fetchURL, conf)
	}
	return healWithAnthropic(site, key, fetchURL, conf)
}

func healWithAnthropic(site, key, fetchURL string, conf LLMConfig) bool {
	client := anthropic.NewClient(option.WithAPIKey(conf.APIKey))
	entry, _ := GetEntry(site, key)
	model := conf.Model
	if model == "" {
		model = string(healerModel)
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(buildPrompt(site, key, entry.XPath, entry.Validate, fetchURL))),
	}
	tools := buildAnthropicTools()

	for round := 0; round < maxRounds; round++ {
		resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
			Model:     anthropic.Model(model),
			MaxTokens: 1024,
			Tools:     tools,
			Messages:  messages,
		})
		if err != nil {
			log.Error(log.Struct{"err", err}, "SelectorHealer: API error")
			return false
		}

		var assistantBlocks []anthropic.ContentBlockParamUnion
		for _, block := range resp.Content {
			if block.Type == "tool_use" {
				assistantBlocks = append(assistantBlocks, anthropic.NewToolUseBlock(block.ID, block.Input, block.Name))
			} else if block.Type == "text" {
				assistantBlocks = append(assistantBlocks, anthropic.NewTextBlock(block.Text))
			}
		}
		messages = append(messages, anthropic.NewAssistantMessage(assistantBlocks...))

		if resp.StopReason == anthropic.StopReasonEndTurn {
			log.Warn(log.Struct{"site", site, "key", key}, "SelectorHealer: agent ended without committing")
			return false
		}

		var toolResults []anthropic.ContentBlockParamUnion
		committed := false
		for _, block := range resp.Content {
			if block.Type != "tool_use" {
				continue
			}
			result, didCommit := dispatchTool(block.Name, block.Input, site, key)
			toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, result, false))
			if didCommit {
				committed = true
			}
		}

		if len(toolResults) > 0 {
			messages = append(messages, anthropic.NewUserMessage(toolResults...))
		}

		if committed {
			log.Info(log.Struct{"site", site, "key", key}, "SelectorHealer: selector repaired")
			return true
		}
	}

	log.Warn(log.Struct{"site", site, "key", key}, "SelectorHealer: exhausted rounds, repair failed")
	return false
}

// oaiMessage mirrors the OpenAI chat message for raw HTTP calls.
type oaiMessage struct {
	Role       string        `json:"role"`
	Content    string        `json:"content,omitempty"`
	ToolCalls  []oaiToolCall `json:"tool_calls,omitempty"`
	ToolCallID string        `json:"tool_call_id,omitempty"`
	Name       string        `json:"name,omitempty"`
}

type oaiToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type oaiResponse struct {
	Choices []struct {
		FinishReason string     `json:"finish_reason"`
		Message      oaiMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func healWithOpenAI(site, key, fetchURL string, conf LLMConfig) bool {
	baseURL := conf.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := conf.Model
	if model == "" {
		model = "gpt-4o-mini"
	}

	entry, _ := GetEntry(site, key)
	messages := []oaiMessage{
		{Role: "user", Content: buildPrompt(site, key, entry.XPath, entry.Validate, fetchURL)},
	}
	tools := buildOpenAIToolsRaw()

	newHTTPClient := func() *http.Client {
		return &http.Client{
			Timeout: 90 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{NextProtos: []string{"http/1.1"}},
				TLSNextProto:    make(map[string]func(string, *tls.Conn) http.RoundTripper),
			},
		}
	}

	for round := 0; round < maxRounds; round++ {
		body, _ := json.Marshal(map[string]any{
			"model":    model,
			"messages": messages,
			"tools":    tools,
		})
		req, _ := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+conf.APIKey)
		req.Header.Set("Content-Type", "application/json")
		log.Debug(log.Struct{"round", round, "bodyLen", len(body)}, "SelectorHealer: sending request")

		resp, err := newHTTPClient().Do(req)
		if err != nil {
			log.Error(log.Struct{"err", err}, "SelectorHealer: API error")
			return false
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var oaiResp oaiResponse
		if err := json.Unmarshal(respBody, &oaiResp); err != nil {
			log.Error(log.Struct{"err", err, "body", string(respBody[:min(len(respBody), 200)])}, "SelectorHealer: parse error")
			return false
		}
		if oaiResp.Error != nil {
			log.Error(log.Struct{"err", oaiResp.Error.Message}, "SelectorHealer: API error")
			return false
		}
		if len(oaiResp.Choices) == 0 {
			log.Warn(log.Struct{"site", site, "key", key}, "SelectorHealer: empty choices")
			return false
		}

		choice := oaiResp.Choices[0]
		messages = append(messages, choice.Message)

		if choice.FinishReason != "tool_calls" {
			log.Warn(log.Struct{"site", site, "key", key}, "SelectorHealer: agent ended without committing")
			return false
		}

		committed := false
		for _, tc := range choice.Message.ToolCalls {
			result, didCommit := dispatchTool(tc.Function.Name, json.RawMessage(tc.Function.Arguments), site, key)
			messages = append(messages, oaiMessage{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
			if didCommit {
				committed = true
			}
		}

		if committed {
			log.Info(log.Struct{"site", site, "key", key}, "SelectorHealer: selector repaired")
			return true
		}
	}

	log.Warn(log.Struct{"site", site, "key", key}, "SelectorHealer: exhausted rounds, repair failed")
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dispatchTool(name string, rawInput json.RawMessage, site, key string) (result string, committed bool) {
	switch name {
	case "fetch_page_snippet":
		var args struct {
			URL  string `json:"url"`
			Hint string `json:"hint"`
		}
		json.Unmarshal(rawInput, &args)
		return toolFetchPageSnippet(args.URL, args.Hint), false

	case "test_xpath":
		var args struct {
			HTML  string `json:"html"`
			XPath string `json:"xpath"`
		}
		json.Unmarshal(rawInput, &args)
		return toolTestXPath(args.HTML, args.XPath), false

	case "validate_result":
		var args struct {
			Value   string `json:"value"`
			Pattern string `json:"pattern"`
		}
		json.Unmarshal(rawInput, &args)
		ok, _ := matchPattern(args.Pattern, args.Value)
		if ok {
			return "VALID", false
		}
		return fmt.Sprintf("INVALID: %q does not match %q", args.Value, args.Pattern), false

	case "commit_selector":
		var args struct {
			XPath string `json:"xpath"`
		}
		json.Unmarshal(rawInput, &args)
		if err := Update(site, key, args.XPath); err != nil {
			return fmt.Sprintf("ERROR: %v", err), false
		}
		return "COMMITTED", true
	}
	return fmt.Sprintf("unknown tool: %s", name), false
}

func toolFetchPageSnippet(url, hint string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("ERROR hint xpath panic %q: %v", hint, r)
		}
	}()
	c := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/112.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 512*1024))

	if hint != "" {
		doc, e := htmlquery.Parse(strings.NewReader(string(body)))
		if e == nil {
			nodes := htmlquery.Find(doc, "//"+hint)
			if len(nodes) > 0 {
				var sb strings.Builder
				for i, n := range nodes {
					if i >= 3 {
						break
					}
					sb.WriteString(htmlquery.OutputHTML(n, true))
				}
				s := sb.String()
				if len(s) > snippetLimit {
					s = s[:snippetLimit]
				}
				return s
			}
		}
	}
	s := string(body)
	if len(s) > snippetLimit {
		s = s[:snippetLimit]
	}
	return s
}

func toolTestXPath(html, xpath string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("ERROR invalid xpath %q: %v", xpath, r)
		}
	}()
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return fmt.Sprintf("ERROR parsing HTML: %v", err)
	}
	nodes := htmlquery.Find(doc, xpath)
	if len(nodes) == 0 {
		return fmt.Sprintf("NO MATCH (0 nodes) for xpath: %q", xpath)
	}
	first := htmlquery.InnerText(nodes[0])
	return fmt.Sprintf("MATCHED %d nodes, first=%q", len(nodes), first)
}

func matchPattern(pattern, value string) (bool, error) {
	if pattern == "" {
		return true, nil
	}
	return regexp.MatchString(pattern, value)
}

func buildPrompt(site, key, oldXPath, validate, fetchURL string) string {
	return fmt.Sprintf(`你是一个网页 XPath 修复专家。

目标网站: %s
需要修复的字段: %s
验证规则 (regex): %s
历史 XPath (已失效): %s
页面 URL: %s

请按步骤操作:
1. 用 fetch_page_snippet 获取页面 HTML（hint 填目标区域的 XPath 片段，如 "table" 或 "ul[@id='infobox']"）
2. 分析 HTML 结构，用 test_xpath 验证候选 XPath
3. 用 validate_result 确认提取值符合验证规则
4. 确认后调用 commit_selector 写入新 XPath
5. 如果找不到有效 XPath，停止并说明原因`,
		site, key, validate, oldXPath, fetchURL)
}

func buildAnthropicTools() []anthropic.ToolUnionParam {
	return []anthropic.ToolUnionParam{
		{OfTool: &anthropic.ToolParam{
			Name:        "fetch_page_snippet",
			Description: anthropic.String("抓取指定 URL 的页面 HTML，返回目标区域片段。hint 填 XPath 片段指定感兴趣的区域（如 table）"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"url":  map[string]string{"type": "string", "description": "要抓取的页面 URL"},
					"hint": map[string]string{"type": "string", "description": "目标区域的 XPath 片段"},
				},
				Required: []string{"url"},
			},
		}},
		{OfTool: &anthropic.ToolParam{
			Name:        "test_xpath",
			Description: anthropic.String("在给定 HTML 上测试 XPath，返回匹配节点数和第一个匹配值"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"html":  map[string]string{"type": "string", "description": "HTML 字符串"},
					"xpath": map[string]string{"type": "string", "description": "要测试的 XPath"},
				},
				Required: []string{"html", "xpath"},
			},
		}},
		{OfTool: &anthropic.ToolParam{
			Name:        "validate_result",
			Description: anthropic.String("用正则验证提取值是否符合预期"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"value":   map[string]string{"type": "string", "description": "待验证的值"},
					"pattern": map[string]string{"type": "string", "description": "Go 正则表达式"},
				},
				Required: []string{"value", "pattern"},
			},
		}},
		{OfTool: &anthropic.ToolParam{
			Name:        "commit_selector",
			Description: anthropic.String("将验证通过的 XPath 写入配置，完成修复"),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"xpath": map[string]string{"type": "string", "description": "验证通过的新 XPath"},
				},
				Required: []string{"xpath"},
			},
		}},
	}
}

func buildOpenAIToolsRaw() []map[string]any {
	defs := []struct {
		name, desc string
		props      map[string]any
		required   []string
	}{
		{
			name: "fetch_page_snippet",
			desc: "抓取指定 URL 的页面 HTML，返回目标区域片段。hint 填 XPath 片段指定感兴趣的区域（如 table）",
			props: map[string]any{
				"url":  map[string]string{"type": "string", "description": "要抓取的页面 URL"},
				"hint": map[string]string{"type": "string", "description": "目标区域的 XPath 片段"},
			},
			required: []string{"url"},
		},
		{
			name: "test_xpath",
			desc: "在给定 HTML 上测试 XPath，返回匹配节点数和第一个匹配值",
			props: map[string]any{
				"html":  map[string]string{"type": "string", "description": "HTML 字符串"},
				"xpath": map[string]string{"type": "string", "description": "要测试的 XPath"},
			},
			required: []string{"html", "xpath"},
		},
		{
			name: "validate_result",
			desc: "用正则验证提取值是否符合预期",
			props: map[string]any{
				"value":   map[string]string{"type": "string", "description": "待验证的值"},
				"pattern": map[string]string{"type": "string", "description": "Go 正则表达式"},
			},
			required: []string{"value", "pattern"},
		},
		{
			name: "commit_selector",
			desc: "将验证通过的 XPath 写入配置，完成修复",
			props: map[string]any{
				"xpath": map[string]string{"type": "string", "description": "验证通过的新 XPath"},
			},
			required: []string{"xpath"},
		},
	}

	var tools []map[string]any
	for _, d := range defs {
		tools = append(tools, map[string]any{
			"type": "function",
			"function": map[string]any{
				"name":        d.name,
				"description": d.desc,
				"parameters": map[string]any{
					"type":       "object",
					"properties": d.props,
					"required":   d.required,
				},
			},
		})
	}
	return tools
}
