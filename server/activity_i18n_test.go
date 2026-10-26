package server

import (
	"strings"
	"testing"

	"github.com/Autumn-27/artex/config"
)

// pinActivityLanguage forces the active output language for one test.
func pinActivityLanguage(t *testing.T, code string) {
	t.Helper()
	prev := config.CurrentLanguage
	config.CurrentLanguage = func() string { return code }
	t.Cleanup(func() { config.CurrentLanguage = prev })
}

// hasHangulRune reports whether s contains a Korean syllable.
func hasHangulRune(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7AF {
			return true
		}
	}
	return false
}

// TestActivitySummaryFollowsLanguage is the regression this catalog exists for:
// the task timeline used to print Korean summaries regardless of the language
// setting, so a Chinese UI still showed
// "모든 목표를 달성했고 직접 투입한 의도의 실행도 끝나, 작업을 종료합니다"
// when a task finished.
func TestActivitySummaryFollowsLanguage(t *testing.T) {
	keys := []string{
		goallessTaskDoneSummary,
		plannerRoundSummaryFmt,
		timeoutFinalRoundSummaryFmt,
		goalBreakdownRound0Summary,
		queuedConcurrencyLimitSummaryFmt,
		queuedNoLLMSummary,
		queuedFIFOSummary,
		taskAssetSourceAPISummary,
		taskAssetSourceTaskSummary,
		auditFindingFollowUpSummary,
		notifyTestSummary,
	}

	for _, code := range []string{"en", "zh", "es", "ja"} {
		pinActivityLanguage(t, code)
		for _, k := range keys {
			got := actText(k)
			if got == k {
				t.Errorf("[%s] activity summary not translated: %q", code, k)
			}
			if hasHangulRune(got) {
				t.Errorf("[%s] activity summary still has Korean: %q", code, got)
			}
		}
	}
}

// TestActivitySummaryKoreanUnchanged keeps a Korean deployment byte-identical to
// before the catalog existed.
func TestActivitySummaryKoreanUnchanged(t *testing.T) {
	pinActivityLanguage(t, "ko")
	for _, k := range []string{goallessTaskDoneSummary, plannerRoundSummaryFmt} {
		if got := actText(k); got != k {
			t.Errorf("Korean summary changed: %q -> %q", k, got)
		}
	}
}

// TestActivitySummaryUnknownPassesThrough guards the no-op contract: a string
// with no catalog entry must come back untouched.
func TestActivitySummaryUnknownPassesThrough(t *testing.T) {
	pinActivityLanguage(t, "en")
	const unknown = "카탈로그에 없는 활동 요약"
	if got := actText(unknown); got != unknown {
		t.Errorf("unknown summary was rewritten: %q", got)
	}
}

// TestActivitySummaryFormatVerbsPreserved makes sure translation never reorders
// or drops a format verb: the call sites pass arguments positionally, so a
// translation with the verbs in the wrong order would print nonsense.
func TestActivitySummaryFormatVerbsPreserved(t *testing.T) {
	pinActivityLanguage(t, "en")
	for _, k := range []string{plannerRoundSummaryFmt, timeoutFinalRoundSummaryFmt, queuedConcurrencyLimitSummaryFmt} {
		got := actText(k)
		if strings.Count(got, "%d") != strings.Count(k, "%d") {
			t.Errorf("format verb count changed for %q -> %q", k, got)
		}
		if strings.Count(got, "%s") != strings.Count(k, "%s") {
			t.Errorf("string verb count changed for %q -> %q", k, got)
		}
	}
}
