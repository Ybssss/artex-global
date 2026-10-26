package server

import (
	"strings"
	"testing"
)

// TestErrCatalogExtraCoversErrorsNew is the regression this file exists for.
//
// These strings are built with errors.New / fmt.Errorf and reach the user through
// err.Error() — written into JSON payloads or into the archive/delivery status
// columns — so writeErr's central trMsg never sees them. Before errCatalogExtra
// they rendered in Korean no matter which language the operator picked.
//
// The test pins the mechanism rather than every string: for a sample that covers
// each surface (archive, notify delivery, side question, intercept, llm pool), the
// Korean source must resolve to something else in every non-Korean language.
func TestErrCatalogExtraCoversErrorsNew(t *testing.T) {
	samples := []string{
		"이 작업은 다른 보관 또는 삭제 작업을 진행하고 있습니다",
		"작업에 아직 실행 중인 에이전트가 있습니다. 먼저 일시정지한 뒤 보관을 다시 시도하세요",
		"보관 스냅샷을 만드는 동안 작업 탐색 기록이 변경되었습니다",
		"Bot Token이 없습니다",
		"Chat ID가 없습니다",
		"Webhook 주소가 없습니다",
		"모델 응답이 중단되었습니다. 다시 질문해 주세요.",
		"모델이 답변을 반환하지 않았습니다.",
		"이미 처리되었거나 존재하지 않는 승인입니다. 기록을 새로 고쳐 주세요",
		"LLM 폴백 체인: 사용 가능한 설정이 없습니다",
		"action 값은 allow, deny, ask 중 하나여야 합니다",
		"한 번에 최대 100개 보관까지 처리할 수 있습니다",
	}

	for _, code := range []string{"en", "zh", "es", "ja"} {
		pinActivityLanguage(t, code)
		for _, ko := range samples {
			got := trMsg(ko)
			if got == ko {
				t.Errorf("[%s] not translated: %q", code, ko)
				continue
			}
			if hasHangulRune(got) {
				t.Errorf("[%s] still Korean: %q -> %q", code, ko, got)
			}
		}
	}
}

// TestErrCatalogExtraKoreanUnchanged keeps a Korean deployment byte-identical.
func TestErrCatalogExtraKoreanUnchanged(t *testing.T) {
	pinActivityLanguage(t, "ko")
	const ko = "Bot Token이 없습니다"
	if got := trMsg(ko); got != ko {
		t.Errorf("Korean message changed: %q -> %q", ko, got)
	}
}

// TestErrCatalogExtraMergedIntoCatalog verifies the init() merge actually ran, so
// trMsg can resolve these strings. Without it the extra map would be dead weight.
func TestErrCatalogExtraMergedIntoCatalog(t *testing.T) {
	if len(errCatalogExtra) == 0 {
		t.Fatal("errCatalogExtra is empty")
	}
	for ko := range errCatalogExtra {
		if _, ok := errCatalog[ko]; !ok {
			t.Errorf("errCatalogExtra key %q was not merged into errCatalog", ko)
		}
	}
}

// TestErrCatalogExtraFormatVerbsPreserved guards the positional-argument contract:
// call sites pass args in the order the Korean template declares them, so a
// translation must keep the same verbs in the same order.
func TestErrCatalogExtraFormatVerbsPreserved(t *testing.T) {
	for ko, langs := range errCatalogExtra {
		wantVerbs := formatVerbs(ko)
		if len(wantVerbs) == 0 {
			continue
		}
		for lang, translated := range langs {
			if translated == "" {
				continue
			}
			if got := formatVerbs(translated); got != wantVerbs {
				t.Errorf("[%s] %q: verbs %q != %q", lang, ko, got, wantVerbs)
			}
		}
	}
}

// formatVerbs returns the ordered %-verbs in s, so two strings can be compared
// for argument-order compatibility.
func formatVerbs(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '%' || i+1 >= len(s) {
			continue
		}
		j := i + 1
		for j < len(s) && strings.ContainsRune("+-# 0123456789.", rune(s[j])) {
			j++
		}
		if j < len(s) {
			b.WriteByte('%')
			b.WriteByte(s[j])
		}
		i = j
	}
	return b.String()
}
