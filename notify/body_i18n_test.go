package notify

import (
	"strings"
	"testing"
)

// pinLang forces the label language for one test and restores the hook after.
func pinLang(t *testing.T, code string) {
	t.Helper()
	prev := Lang
	Lang = func() string { return code }
	t.Cleanup(func() { Lang = prev })
}

// TestBodyTextFollowsLanguage is the regression this catalog exists for: the
// digest heading, the footer and the per-item field labels used to stay Korean
// no matter which language the channel was configured for.
func TestBodyTextFollowsLanguage(t *testing.T) {
	keys := []string{
		"취약점 요약 · 총 %d건",
		"취약점 알림",
		"**신규 취약점 %d건**",
		"**유형**: %s\n",
		"**자산**: %s\n",
		"**개요**: %s\n",
		"[상세 보기](%s)\n",
	}
	for _, code := range []string{"en", "zh", "es", "ja"} {
		pinLang(t, code)
		for _, k := range keys {
			got := bodyText(k)
			if got == k {
				t.Errorf("[%s] %q not translated", code, k)
			}
			if hasHangul(got) {
				t.Errorf("[%s] %q still has Korean: %q", code, k, got)
			}
		}
	}
}

// TestBodyTextKoreanUnchanged keeps a Korean deployment byte-identical.
func TestBodyTextKoreanUnchanged(t *testing.T) {
	pinLang(t, "ko")
	for _, k := range []string{"취약점 알림", "**유형**: %s\n"} {
		if got := bodyText(k); got != k {
			t.Errorf("Korean output changed: %q -> %q", k, got)
		}
	}
}

// TestBodyTextUnknownKeyPassesThrough guards the no-op contract: a string with no
// catalog entry must come back untouched rather than becoming empty.
func TestBodyTextUnknownKeyPassesThrough(t *testing.T) {
	pinLang(t, "en")
	const unknown = "카탈로그에 없는 문자열"
	if got := bodyText(unknown); got != unknown {
		t.Errorf("unknown key was rewritten: %q", got)
	}
}

// TestMarkdownBodyLocalized renders a real digest body in English and asserts no
// Korean survives anywhere in it.
func TestMarkdownBodyLocalized(t *testing.T) {
	pinLang(t, "en")
	m := Message{
		HomeURL:       "https://example.test",
		Batch:         true,
		WindowMinutes: 30,
		Items: []Item{{
			Severity:  "high",
			VulnClass: "SQL Injection",
			Summary:   "a finding",
			DetailURL: "https://example.test/f/1",
		}},
	}
	body, _ := markdownBody(m, 4096)
	if hasHangul(body) {
		t.Errorf("markdown body still contains Korean:\n%s", body)
	}
	if !strings.Contains(strings.ToLower(body), "vulnerab") {
		t.Errorf("digest heading not translated:\n%s", body)
	}
}

// TestHTMLBodyLocalized is the same contract for the HTML channel.
func TestHTMLBodyLocalized(t *testing.T) {
	pinLang(t, "zh")
	m := Message{
		HomeURL:       "https://example.test",
		Batch:         true,
		WindowMinutes: 30,
		Items: []Item{{
			Severity:  "high",
			VulnClass: "SQL 注入",
			Summary:   "一个发现",
			DetailURL: "https://example.test/f/1",
		}},
	}
	body := htmlBody(m, 4096)
	if hasHangul(body) {
		t.Errorf("html body still contains Korean:\n%s", body)
	}
}
