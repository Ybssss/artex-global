package agent

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// bodyOf drains a rewritten request body for assertions.
func bodyOf(t *testing.T, req *http.Request) string {
	t.Helper()
	rc, err := req.GetBody()
	if err != nil {
		t.Fatalf("GetBody: %v", err)
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

// messagesOf decodes the messages array out of a rewritten body.
func messagesOf(t *testing.T, body string) []map[string]any {
	t.Helper()
	var doc struct {
		Messages []map[string]any `json:"messages"`
	}
	if err := json.Unmarshal([]byte(body), &doc); err != nil {
		t.Fatalf("decode body: %v\nbody=%s", err, body)
	}
	return doc.Messages
}

// chatRequest builds a POST to /chat/completions carrying body, with GetBody
// populated the way norma's *bytes.Reader produces it.
func chatRequest(body string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "https://gateway.test/v1/chat/completions", strings.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(body)), nil
	}
	return req
}

// TestToolRoleFoldRewritesToolMessages is the regression this whole file exists
// for: a `role:"tool"` message must reach the gateway as user text, because some
// gateways silently drop the native shape and the model then reports that its
// tools returned nothing.
func TestToolRoleFoldRewritesToolMessages(t *testing.T) {
	t.Setenv("ARTEX_LLM_FOLD_TOOL_ROLE", "1")

	in := `{"model":"m","messages":[` +
		`{"role":"user","content":"go"},` +
		`{"role":"assistant","content":"","tool_calls":[{"id":"call_X","type":"function","function":{"name":"Bash","arguments":"{}"}}]},` +
		`{"role":"tool","tool_call_id":"call_X","content":"CANARY-PAIR-5W8"},` +
		`{"role":"user","content":"repeat"}]}`

	var seen string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		seen = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	// Build a real client request (not httptest.NewRequest, which sets RequestURI
	// and cannot be used by a client) so the transport sees what norma produces.
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/v1/chat/completions", strings.NewReader(in))
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader(in)), nil
	}
	client := &http.Client{Transport: openAIToolRoleNormalizer{
		base:    http.DefaultTransport,
		enabled: func() bool { return true },
	}}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	_ = resp.Body.Close()

	if strings.Contains(seen, `"role":"tool"`) {
		t.Fatalf("tool role survived the rewrite: %s", seen)
	}
	msgs := messagesOf(t, seen)
	if len(msgs) != 4 {
		t.Fatalf("want 4 messages, got %d: %s", len(msgs), seen)
	}
	if msgs[2]["role"] != "user" {
		t.Fatalf("tool message became %v, want user", msgs[2]["role"])
	}
	content, _ := msgs[2]["content"].(string)
	if !strings.Contains(content, "CANARY-PAIR-5W8") {
		t.Fatalf("canary lost: %q", content)
	}
	if !strings.Contains(content, "Bash") {
		t.Fatalf("tool name not carried: %q", content)
	}
	// The assistant turn must keep its tool_calls so the transcript still records
	// which call produced the output.
	if _, ok := msgs[1]["tool_calls"]; !ok {
		t.Fatalf("assistant tool_calls dropped: %v", msgs[1])
	}
}

// TestToolRoleFoldOffByDefault keeps native tool messages intact on a correct
// gateway. Turning this on unconditionally would lose call-id correlation for
// every caller that never needed the workaround.
func TestToolRoleFoldOffByDefault(t *testing.T) {
	t.Setenv("ARTEX_LLM_FOLD_TOOL_ROLE", "")

	in := `{"model":"m","messages":[` +
		`{"role":"tool","tool_call_id":"c1","content":"X"}]}`
	req := chatRequest(in)
	if got := bodyOf(t, req); got != in {
		t.Fatalf("disabled normalizer changed the body: %s", got)
	}
	if toolRoleFoldEnabled() {
		t.Fatal("fold must be off when the env var is unset")
	}
}

// TestToolRoleFoldLeavesOtherRequestsAlone guards the no-op contract: a request
// with no tool message, a non-chat path, and a non-JSON body must all pass
// through byte for byte.
func TestToolRoleFoldLeavesOtherRequestsAlone(t *testing.T) {
	plain := `{"model":"m","messages":[{"role":"user","content":"hi"}]}`
	if out, changed := foldToolResultsIntoUserTurns([]byte(plain)); changed || string(out) != plain {
		t.Fatalf("plain body changed: %s", out)
	}
	if _, changed := foldToolResultsIntoUserTurns([]byte("not json")); changed {
		t.Fatal("non-JSON body reported a change")
	}
	if isChatCompletionsPath("/v1/messages") {
		t.Fatal("anthropic path must not match")
	}
	if !isChatCompletionsPath("/v1/chat/completions") {
		t.Fatal("chat path must match")
	}
}

// TestToolRoleFoldHandlesContentParts covers the other OpenAI spelling, where a
// tool message's content is an array of text parts.
func TestToolRoleFoldHandlesContentParts(t *testing.T) {
	in := `{"model":"m","messages":[` +
		`{"role":"tool","tool_call_id":"c1","content":[{"type":"text","text":"PART-ONE"},{"type":"text","text":"-PART-TWO"}]}]}`
	out, changed := foldToolResultsIntoUserTurns([]byte(in))
	if !changed {
		t.Fatal("content-parts tool message was not folded")
	}
	msgs := messagesOf(t, string(out))
	content, _ := msgs[0]["content"].(string)
	if !strings.Contains(content, "PART-ONE-PART-TWO") {
		t.Fatalf("parts not concatenated: %q", content)
	}
}

// TestToolRoleFoldKeepsEmptyResultsHonest: an empty tool output must still be
// represented. Silence is indistinguishable from a dropped message, which is the
// exact failure mode this workaround exists to prevent.
func TestToolRoleFoldKeepsEmptyResultsHonest(t *testing.T) {
	in := `{"model":"m","messages":[{"role":"tool","tool_call_id":"c1","content":""}]}`
	out, changed := foldToolResultsIntoUserTurns([]byte(in))
	if !changed {
		t.Fatal("empty tool result was dropped")
	}
	msgs := messagesOf(t, string(out))
	content, _ := msgs[0]["content"].(string)
	if !strings.Contains(content, "(no output)") {
		t.Fatalf("empty result not marked: %q", content)
	}
}

// TestToolRoleFoldBindsParallelCalls keeps two tool results in order and
// attributed to their own calls.
func TestToolRoleFoldBindsParallelCalls(t *testing.T) {
	in := `{"model":"m","messages":[` +
		`{"role":"assistant","content":"","tool_calls":[` +
		`{"id":"c1","type":"function","function":{"name":"Bash","arguments":"{}"}},` +
		`{"id":"c2","type":"function","function":{"name":"Read","arguments":"{}"}}]},` +
		`{"role":"tool","tool_call_id":"c1","content":"OUT-ONE"},` +
		`{"role":"tool","tool_call_id":"c2","content":"OUT-TWO"}]}`
	out, changed := foldToolResultsIntoUserTurns([]byte(in))
	if !changed {
		t.Fatal("parallel tool results were not folded")
	}
	msgs := messagesOf(t, string(out))
	if len(msgs) != 3 {
		t.Fatalf("want 3 messages, got %d", len(msgs))
	}
	one, _ := msgs[1]["content"].(string)
	two, _ := msgs[2]["content"].(string)
	if !strings.Contains(one, "OUT-ONE") || !strings.Contains(one, "Bash") {
		t.Fatalf("first result misattributed: %q", one)
	}
	if !strings.Contains(two, "OUT-TWO") || !strings.Contains(two, "Read") {
		t.Fatalf("second result misattributed: %q", two)
	}
}
