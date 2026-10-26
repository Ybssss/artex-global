package report

import (
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/config"
)

// pinLanguage forces the active output language for one test and restores the
// previous hook afterwards.
func pinLanguage(t *testing.T, code string) {
	t.Helper()
	prev := config.CurrentLanguage
	config.CurrentLanguage = func() string { return code }
	t.Cleanup(func() { config.CurrentLanguage = prev })
}

// hasHangul reports whether s contains any Korean syllable.
func hasHangul(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7AF {
			return true
		}
	}
	return false
}

// TestMarkdownFollowsActiveLanguage is the regression this package's i18n exists
// for: a report rendered under a non-Korean language must not contain Korean.
// Before the fix every literal was hardcoded, so this failed for en/zh/es/ja.
func TestMarkdownFollowsActiveLanguage(t *testing.T) {
	in := Input{
		Title:       "Demo task",
		Goal:        "check everything",
		GeneratedAt: time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC),
		AssetCounts: map[string]int{"ip": 2},
	}

	for _, code := range []string{"en", "zh", "es", "ja"} {
		pinLanguage(t, code)
		out := Markdown(in)
		if hasHangul(out) {
			t.Errorf("[%s] report still contains Korean:\n%s", code, out)
		}
	}
}

// TestMarkdownKoreanUnchanged pins the other half of the contract: a Korean
// deployment must render byte-identically to before the catalog existed.
func TestMarkdownKoreanUnchanged(t *testing.T) {
	pinLanguage(t, "ko")
	out := Markdown(Input{Title: "데모 작업", GeneratedAt: time.Now()})
	for _, want := range []string{"# 침투 테스트 보고서: 데모 작업", "## 요약", "## 취약점"} {
		if !strings.Contains(out, want) {
			t.Errorf("Korean output lost %q:\n%s", want, out)
		}
	}
}

// TestMarkdownLocalizedHeadings checks that each language actually renders its
// own heading rather than falling back to English.
func TestMarkdownLocalizedHeadings(t *testing.T) {
	in := Input{Title: "T", GeneratedAt: time.Now()}

	cases := map[string]string{
		"en": "# Penetration Test Report: T",
		"zh": "# 渗透测试报告：T",
		"es": "# Informe de prueba de penetración: T",
		"ja": "# ペネトレーションテストレポート: T",
	}
	for code, want := range cases {
		pinLanguage(t, code)
		if out := Markdown(in); !strings.Contains(out, want) {
			t.Errorf("[%s] missing %q, got:\n%s", code, want, out)
		}
	}
}

// TestCSVHeaderFollowsLanguage covers the export path, whose header was also
// hardcoded Korean.
func TestCSVHeaderFollowsLanguage(t *testing.T) {
	pinLanguage(t, "en")
	csv := string(FindingsCSV(nil))
	if hasHangul(csv) {
		t.Errorf("CSV header still contains Korean:\n%s", csv)
	}
	if !strings.Contains(csv, "Severity") {
		t.Errorf("CSV header not translated:\n%s", csv)
	}
}
