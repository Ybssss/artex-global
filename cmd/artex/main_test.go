package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/Autumn-27/artex/agent"
)

func TestShutdownContextPreservesNamedCause(t *testing.T) {
	signalCtx, signalCancel := context.WithCancel(context.Background())
	ctx, shutdown := shutdownContext(signalCtx)
	defer shutdown(nil)
	signalCancel()
	<-ctx.Done()
	if code, _, _, ok := agent.AbortReason(ctx); !ok || code != "shutdown" {
		t.Fatalf("code=%q ok=%v, want shutdown", code, ok)
	}
}

// TestPrintBannerFollowsLanguage verifies the startup banner is rendered in the
// language the operator configured, rather than always in Korean.
//
// The banner used to be hardcoded Korean, and this test used to assert exactly
// that. Requiring Korean made the language setting a lie: an operator running
// with ARTEX_LANG=en still saw Korean as the first line of output. What matters
// is that the banner tracks config.Language(), so the test pins both ends:
// Korean under ko, and no Korean at all under the other four.
func TestPrintBannerFollowsLanguage(t *testing.T) {
	t.Setenv("ARTEX_LANG", "ko")
	ko := captureStdout(t, func() { printBanner(":8787") })
	if !strings.ContainsFunc(ko, func(r rune) bool { return unicode.Is(unicode.Hangul, r) }) {
		t.Fatalf("banner is not Korean under ARTEX_LANG=ko: %q", ko)
	}
	for _, want := range []string{"자율 침투 테스트 시스템", "버전", "수신 대기", ":8787"} {
		if !strings.Contains(ko, want) {
			t.Errorf("Korean banner missing %q: %q", want, ko)
		}
	}
	// Even in Korean, no CJK Han residue: the brief requires pure Hangul here.
	for _, r := range ko {
		if unicode.Is(unicode.Han, r) {
			t.Errorf("Korean banner has CJK Han residue (%q): %q", r, ko)
		}
	}

	for _, code := range []string{"en", "zh", "es", "ja"} {
		t.Setenv("ARTEX_LANG", code)
		out := captureStdout(t, func() { printBanner(":8787") })
		if strings.ContainsFunc(out, func(r rune) bool { return unicode.Is(unicode.Hangul, r) }) {
			t.Errorf("[%s] banner still contains Korean: %q", code, out)
		}
		if !strings.Contains(out, ":8787") {
			t.Errorf("[%s] banner dropped the listen address: %q", code, out)
		}
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns whatever
// fn wrote. printBanner uses fmt.Print*, which resolves os.Stdout at call time,
// so swapping it here captures the banner.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	os.Stdout = orig
	return <-done
}
