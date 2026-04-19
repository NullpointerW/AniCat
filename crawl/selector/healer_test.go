package selector

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/NullpointerW/anicat/log"
)

func TestMain(m *testing.M) {
	log.Init("text", "debug", "2006-01-02T15:04:05", false, os.Stdout)
	os.Exit(m.Run())
}

// ── selectors.yaml 完整性校验 ────────────────────────────────────────

func TestSelectorsYamlCoverage(t *testing.T) {
	// 找到项目根目录的 selectors.yaml
	yamlPath := "../../selectors.yaml"
	if err := Load(yamlPath); err != nil {
		t.Fatalf("load selectors.yaml: %v", err)
	}

	required := []struct {
		site, key    string
		sampleValid  string // 用 validate 规则验证的合法样本（空表示跳过）
		sampleInvalid string // 应该失败的样本（空表示跳过）
	}{
		// mikan
		{"mikan", "rss_li", "", ""},
		{"mikan", "rss_endpoint", "/RSS/Bangumi?bangumiId=1&subgroupid=2", "not-a-path"},
		{"mikan", "subgroup_container", "", ""},
		{"mikan", "bgm_url", "https://bgm.tv/subject/123", "https://example.com"},
		{"mikan", "search_rows", "", ""},
		{"mikan", "magnet_link", "magnet:?xt=urn:btih:abc", "http://example.com"},
		{"mikan", "filename", "", ""},
		{"mikan", "filesize", "781.5MB", "2024-01-01"},
		{"mikan", "update_time", "", ""},
		{"mikan", "rss_table_rows", "", ""},
		{"mikan", "rss_item_name", "", ""},
		{"mikan", "rss_item_size", "643.2MB", "2024-01-01"},
		{"mikan", "rss_item_uptime", "", ""},
		// bgmtv
		{"bgmtv", "infobox", "话数", ""},
		{"bgmtv", "origin_name", "Frieren", ""},
		// tmdb
		{"tmdb", "title", "葬送的芙莉莲", ""},
		{"tmdb", "release_date", "2023-09-29", ""},
	}

	for _, tc := range required {
		xpath := Get(tc.site, tc.key)
		if xpath == "" {
			t.Errorf("[%s/%s] XPath 为空，selectors.yaml 缺少该 key", tc.site, tc.key)
			continue
		}
		t.Logf("[%s/%s] xpath=%q", tc.site, tc.key, xpath)

		if tc.sampleValid != "" {
			if !Validate(tc.site, tc.key, tc.sampleValid) {
				t.Errorf("[%s/%s] 合法样本 %q 未通过 validate 规则", tc.site, tc.key, tc.sampleValid)
			}
		}
		if tc.sampleInvalid != "" {
			if Validate(tc.site, tc.key, tc.sampleInvalid) {
				t.Errorf("[%s/%s] 非法样本 %q 意外通过了 validate 规则", tc.site, tc.key, tc.sampleInvalid)
			}
		}
	}
}

// ── 纯逻辑：matchPattern ──────────────────────────────────────────────

func TestMatchPattern(t *testing.T) {
	cases := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"MB|GB", "642.5MB", true},
		{"MB|GB", "34.2 GB", true},
		{"MB|GB", "2006 年 04 月", false},  // 日期误入 size 字段
		{"^magnet:\\?", "magnet:?xt=urn:btih:abc", true},
		{"^magnet:\\?", "http://example.com", false},
		{"", "anything", true}, // 空规则永远通过
	}
	for _, c := range cases {
		got, err := matchPattern(c.pattern, c.value)
		if err != nil {
			t.Errorf("pattern=%q value=%q unexpected err: %v", c.pattern, c.value, err)
		}
		if got != c.want {
			t.Errorf("pattern=%q value=%q: got %v, want %v", c.pattern, c.value, got, c.want)
		}
	}
}

// ── 工具函数：toolTestXPath ───────────────────────────────────────────

func TestToolTestXPath_Match(t *testing.T) {
	html := `<html><body>
		<table><tbody>
			<tr><td>checkbox</td><td>番剧名称.mkv</td><td>642.5MB</td><td>2024-01-01</td></tr>
		</tbody></table>
	</body></html>`

	result := toolTestXPath(html, `//table/tbody/tr/td[3]`)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	t.Log(result)
	// 应该包含 MB
	if !contains(result, "642.5MB") {
		t.Errorf("expected result to contain '642.5MB', got: %s", result)
	}
}

func TestToolTestXPath_NoMatch(t *testing.T) {
	html := `<html><body><p>hello</p></body></html>`
	result := toolTestXPath(html, `//table/tbody/tr/td[3]`)
	if !contains(result, "NO MATCH") {
		t.Errorf("expected NO MATCH, got: %s", result)
	}
}

func TestToolTestXPath_InvalidXPath(t *testing.T) {
	result := toolTestXPath("<html></html>", `///invalid[[[`)
	if !contains(result, "NO MATCH") && !contains(result, "ERROR") {
		t.Errorf("expected error or no match for invalid XPath, got: %s", result)
	}
}

// ── 工具函数：toolFetchPageSnippet（需要网络）────────────────────────

func TestToolFetchPageSnippet_WithHint(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	result := toolFetchPageSnippet("https://mikanime.tv/Home/Search?searchstr=%E8%91%AC%E9%80%81%E7%9A%84%E8%8A%99%E8%8E%89%E8%8E%B2", "table")
	if result == "" {
		t.Fatal("expected non-empty snippet")
	}
	t.Logf("snippet length=%d, preview=%.200s", len(result), result)
	if !contains(result, "<table") && !contains(result, "ERROR") {
		t.Errorf("expected HTML table snippet or error, got: %.100s", result)
	}
}

// ── selector 加载与热更新 ────────────────────────────────────────────

func TestLoadAndGet(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	yaml := `
mikan:
  filesize:
    xpath: "//td[3]"
    validate: "MB|GB"
`
	if err := os.WriteFile(p, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Load(p); err != nil {
		t.Fatal(err)
	}
	got := Get("mikan", "filesize")
	if got != "//td[3]" {
		t.Errorf("Get: got %q, want %q", got, "//td[3]")
	}
}

func TestUpdate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte("mikan:\n  filesize:\n    xpath: \"//td[3]\"\n    validate: \"MB|GB\"\n"), 0644)
	Load(p)

	newXPath := "//table/tbody/tr/td[3]"
	if err := Update("mikan", "filesize", newXPath); err != nil {
		t.Fatal(err)
	}
	if got := Get("mikan", "filesize"); got != newXPath {
		t.Errorf("after Update: got %q, want %q", got, newXPath)
	}
	// verify persisted to disk
	Load(p)
	if got := Get("mikan", "filesize"); got != newXPath {
		t.Errorf("after reload: got %q, want %q", got, newXPath)
	}
}

func TestValidate(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte("mikan:\n  filesize:\n    xpath: \"//td[3]\"\n    validate: \"MB|GB\"\n"), 0644)
	Load(p)

	if !Validate("mikan", "filesize", "642.5MB") {
		t.Error("expected 642.5MB to pass validation")
	}
	if Validate("mikan", "filesize", "2024-01-01") {
		t.Error("expected date string to fail validation")
	}
}

// ── 完整 agent loop（需要 ANTHROPIC_API_KEY）────────────────────────

func TestHealMikanFilesize(t *testing.T) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping integration test")
	}
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// 故意写一个错误的 XPath（td[1] 拿到的是文件名，不含 MB/GB）
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte(`mikan:
  filesize:
    xpath: "//table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]"
    validate: "MB|GB"
`), 0644)
	Load(p)

	searchURL := "https://mikanime.tv/Home/Search?searchstr=%E8%91%AC%E9%80%81%E7%9A%84%E8%8A%99%E8%8E%89%E8%8E%B2"
	healed := Heal("mikan", "filesize", searchURL, LLMConfig{Style: "anthropic", APIKey: apiKey})
	if !healed {
		t.Error("Heal returned false — agent failed to repair the selector")
		return
	}

	newXPath := Get("mikan", "filesize")
	t.Logf("new xpath: %s", newXPath)

	// 验证修复后的 XPath 确实能抓到含 MB/GB 的值
	result := toolTestXPath(toolFetchPageSnippet(searchURL, "table"), newXPath)
	t.Logf("test result: %s", result)
	if !contains(result, "MB") && !contains(result, "GB") {
		t.Errorf("repaired xpath still not matching size: %s", result)
	}
}

func TestHealMikanFilesizeQwen(t *testing.T) {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		t.Skip("QWEN_API_KEY not set, skipping integration test")
	}
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte(`mikan:
  filesize:
    xpath: "//table[@class='table table-striped tbl-border fadeIn']/tbody/tr[@class='js-search-results-row'][1]/td[1]"
    validate: "MB|GB"
`), 0644)
	Load(p)

	searchURL := "https://mikanime.tv/Home/Search?searchstr=%E8%91%AC%E9%80%81%E7%9A%84%E8%8A%99%E8%8E%89%E8%8E%B2"
	conf := LLMConfig{
		Style:   "openai",
		APIKey:  apiKey,
		Model:   os.Getenv("QWEN_MODEL"),
		BaseURL: os.Getenv("QWEN_BASE_URL"),
	}
	healed := Heal("mikan", "filesize", searchURL, conf)
	if !healed {
		t.Error("Heal returned false — agent failed to repair the selector")
		return
	}

	newXPath := Get("mikan", "filesize")
	t.Logf("new xpath: %s", newXPath)

	result := toolTestXPath(toolFetchPageSnippet(searchURL, "table"), newXPath)
	t.Logf("test result: %s", result)
	if !contains(result, "MB") && !contains(result, "GB") {
		t.Errorf("repaired xpath still not matching size: %s", result)
	}
}

// TestHealWithMockServer 使用本地 mock 服务器验证完整修复流程：
//  1. 提供一个 filesize 列位置变化的 HTML（td[4] 而非原来的 td[3]）
//  2. 加载故意错误的 XPath，验证爬取失败
//  3. 调用 LLM Healer 自动修复
//  4. 再次爬取验证成功
func TestHealWithMockServer(t *testing.T) {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		t.Skip("QWEN_API_KEY not set, skipping integration test")
	}
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// ── 1. 启动 mock 服务器，filesize 在 td[4]（日期移到了 td[3]）──────
	mockHTML := `<!DOCTYPE html><html><body>
<table class="table table-striped tbl-border fadeIn">
  <tbody>
    <tr class="js-search-results-row">
      <td><input type="checkbox"/></td>
      <td><a class="magnet-link-wrap">葬送的芙莉莲 01.mkv</a></td>
      <td>2024-01-15</td>
      <td>781.5MB</td>
    </tr>
    <tr class="js-search-results-row">
      <td><input type="checkbox"/></td>
      <td><a class="magnet-link-wrap">葬送的芙莉莲 02.mkv</a></td>
      <td>2024-01-22</td>
      <td>643.2MB</td>
    </tr>
  </tbody>
</table>
</body></html>`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, mockHTML)
	}))
	defer ts.Close()

	// ── 2. 加载错误的 XPath：td[3] 拿到的是日期，不含 MB/GB ──────────
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte(`mikan:
  filesize:
    xpath: "//table/tbody/tr[1]/td[3]"
    validate: "MB|GB"
`), 0644)
	Load(p)

	// ── 3. 第一次爬取：验证当前 XPath 确实拿到错误数据 ──────────────
	pageHTML := toolFetchPageSnippet(ts.URL, "table")
	if contains(pageHTML, "ERROR") {
		t.Fatalf("mock server fetch failed: %s", pageHTML)
	}
	before := toolTestXPath(pageHTML, Get("mikan", "filesize"))
	t.Logf("修复前 xpath=%q  结果: %s", Get("mikan", "filesize"), before)
	if contains(before, "MB") || contains(before, "GB") {
		t.Fatal("期望修复前 XPath 拿不到 MB/GB，但却拿到了，HTML 结构可能不对")
	}

	// ── 4. 调用 LLM Healer 修复 ──────────────────────────────────────
	conf := LLMConfig{
		Style:   "openai",
		APIKey:  apiKey,
		Model:   os.Getenv("QWEN_MODEL"),
		BaseURL: os.Getenv("QWEN_BASE_URL"),
	}
	healed := Heal("mikan", "filesize", ts.URL, conf)
	if !healed {
		t.Fatal("Heal 返回 false，LLM 未能修复选择器")
	}
	newXPath := Get("mikan", "filesize")
	t.Logf("修复后 xpath: %s", newXPath)

	// ── 5. 第二次爬取：验证新 XPath 能正确拿到文件大小 ──────────────
	after := toolTestXPath(pageHTML, newXPath)
	t.Logf("修复后爬取结果: %s", after)
	if !contains(after, "MB") && !contains(after, "GB") {
		t.Errorf("修复后 XPath 仍无法匹配到 MB/GB: %s", after)
	}
}

// TestHealWithMockRssServer 使用 RSS 列表页结构验证修复流程：
//  1. Mock HTML 中 filesize 列移到了 td[4]（原为 td[3]）
//  2. 错误 XPath 拿到日期，不含 MB/GB
//  3. Healer 修复后再次爬取成功
func TestHealWithMockRssServer(t *testing.T) {
	apiKey := os.Getenv("QWEN_API_KEY")
	if apiKey == "" {
		t.Skip("QWEN_API_KEY not set, skipping integration test")
	}
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// RSS 列表页 mock HTML：filesize 在 td[4]，日期在 td[3]
	mockHTML := `<!DOCTYPE html>
<html><body class="main">
<div id="sk-container">
<div class="central-container">
  <div class="subgroup-text">
    <a>喵萌奶茶屋</a>
    <a class="mikan-rss" href="/RSS/Bangumi?bangumiId=3437&subgroupid=370">RSS</a>
    <table class="table table-striped tbl-border fadeIn">
      <tbody>
        <tr>
          <td><input type="checkbox"/></td>
          <td><a class="magnet-link-wrap">葬送的芙莉莲 01.mkv</a></td>
          <td>2024-01-15</td>
          <td>781.5MB</td>
        </tr>
        <tr>
          <td><input type="checkbox"/></td>
          <td><a class="magnet-link-wrap">葬送的芙莉莲 02.mkv</a></td>
          <td>2024-01-22</td>
          <td>643.2MB</td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
</div>
</body></html>`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, mockHTML)
	}))
	defer ts.Close()

	// 错误的 XPath：td[3] 拿到日期
	dir := t.TempDir()
	p := filepath.Join(dir, "selectors.yaml")
	os.WriteFile(p, []byte(`mikan:
  rss_item_size:
    xpath: "/td[3]"
    validate: "MB|GB"
`), 0644)
	Load(p)

	// 第一次爬取：验证错误 XPath 拿到日期
	pageHTML := toolFetchPageSnippet(ts.URL, "table")
	if contains(pageHTML, "ERROR") {
		t.Fatalf("mock server fetch failed: %s", pageHTML)
	}
	before := toolTestXPath(pageHTML, Get("mikan", "rss_item_size"))
	t.Logf("修复前 xpath=%q  结果: %s", Get("mikan", "rss_item_size"), before)
	if contains(before, "MB") || contains(before, "GB") {
		t.Fatal("期望修复前拿不到 MB/GB，但却拿到了")
	}

	// 调用 Healer
	conf := LLMConfig{
		Style:   "openai",
		APIKey:  apiKey,
		Model:   os.Getenv("QWEN_MODEL"),
		BaseURL: os.Getenv("QWEN_BASE_URL"),
	}
	healed := Heal("mikan", "rss_item_size", ts.URL, conf)
	if !healed {
		t.Fatal("Heal 返回 false，LLM 未能修复选择器")
	}
	newXPath := Get("mikan", "rss_item_size")
	t.Logf("修复后 xpath: %s", newXPath)

	// 第二次爬取：验证修复后能拿到文件大小
	after := toolTestXPath(pageHTML, newXPath)
	t.Logf("修复后爬取结果: %s", after)
	if !contains(after, "MB") && !contains(after, "GB") {
		t.Errorf("修复后 XPath 仍无法匹配到 MB/GB: %s", after)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}

