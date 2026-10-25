package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

// openAIToolRoleNormalizer rewrites outgoing OpenAI-format chat requests so that
// tool results travel as ordinary user turns instead of `role:"tool"` messages.
//
// Why this exists: some OpenAI-compatible gateways drop `role:"tool"` messages
// while assembling the upstream prompt. The request still succeeds, the tool
// content simply never reaches the model, and the model then reports that its
// tools returned nothing. A caller cannot distinguish that from a broken
// toolchain, so an agent loop can spin for dozens of turns chasing a fault that
// is not on its side. Folding the result into a user turn is the one shape every
// gateway has to carry, because it is just text.
//
// What is preserved: the assistant message keeps its `tool_calls`, so the
// conversation still records which call produced which output. What is lost: the
// `tool_call_id` correlation on the wire. That is acceptable here because the
// agent reads tool output as text and never replays a tool result to the
// provider — but it is why this rewrite is opt-in per profile rather than always
// on.
//
// The rewrite is a no-op for anything it does not recognise: a body that is not
// JSON, a non-chat endpoint, or a request with no `role:"tool"` message is
// returned untouched, byte for byte.
type openAIToolRoleNormalizer struct {
	base http.RoundTripper
	// enabled is read per request so a profile change takes effect without
	// rebuilding the client.
	enabled func() bool
}

func (t openAIToolRoleNormalizer) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.enabled == nil || !t.enabled() || req.GetBody == nil {
		return t.base.RoundTrip(req)
	}
	if !isChatCompletionsPath(req.URL.Path) {
		return t.base.RoundTrip(req)
	}
	original, err := readRequestBody(req)
	if err != nil {
		return t.base.RoundTrip(req)
	}
	rewritten, changed := foldToolResultsIntoUserTurns(original)
	if !changed {
		return t.base.RoundTrip(req)
	}
	req.Body = io.NopCloser(bytes.NewReader(rewritten))
	req.ContentLength = int64(len(rewritten))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(rewritten)), nil
	}
	req.Header.Set("Content-Length", itoa(len(rewritten)))
	return t.base.RoundTrip(req)
}

func isChatCompletionsPath(path string) bool {
	return strings.HasSuffix(strings.TrimRight(path, "/"), "/chat/completions")
}

// toolRoleFoldEnabled reports whether outgoing chat requests should have their
// tool results folded into user turns. Off by default: it is a workaround for a
// specific class of gateway, and a caller on a correct gateway should keep the
// native `role:"tool"` shape (including its call-id correlation).
//
// Enabled per task/profile via ARTEX_LLM_FOLD_TOOL_ROLE=1|true|on|yes. Read per
// request so flipping the env var does not require a rebuild.
func toolRoleFoldEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ARTEX_LLM_FOLD_TOOL_ROLE"))) {
	case "1", "true", "on", "yes":
		return true
	}
	return false
}

// readRequestBody returns the body bytes without consuming the request: the
// snapshot is taken from GetBody so the caller's Body stays readable.
func readRequestBody(req *http.Request) ([]byte, error) {
	rc, err := req.GetBody()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// foldToolResultsIntoUserTurns rewrites every `role:"tool"` message into a
// `role:"user"` message carrying the same text, prefixed with the tool name when
// it can be recovered from the preceding assistant `tool_calls`. It reports
// whether anything changed.
//
// The result is shaped as a user turn rather than merged into a neighbouring
// user turn on purpose: keeping message boundaries intact means the model sees
// the same turn structure it would have seen with a real tool result, and it
// avoids reordering content relative to an assistant message.
func foldToolResultsIntoUserTurns(body []byte) ([]byte, bool) {
	var doc map[string]json.RawMessage
	if err := json.Unmarshal(body, &doc); err != nil {
		return body, false
	}
	rawMessages, ok := doc["messages"]
	if !ok {
		return body, false
	}
	var messages []json.RawMessage
	if err := json.Unmarshal(rawMessages, &messages); err != nil {
		return body, false
	}

	// toolNames maps a call id to the function name that issued it, so a folded
	// result can still say which tool produced it.
	toolNames := map[string]string{}
	changed := false
	out := make([]json.RawMessage, 0, len(messages))

	for _, raw := range messages {
		var msg struct {
			Role       string          `json:"role"`
			Content    json.RawMessage `json:"content"`
			ToolCallID string          `json:"tool_call_id"`
			ToolCalls  []struct {
				ID       string `json:"id"`
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
			} `json:"tool_calls"`
		}
		if err := json.Unmarshal(raw, &msg); err != nil {
			out = append(out, raw)
			continue
		}
		if msg.Role == "assistant" {
			for _, call := range msg.ToolCalls {
				if call.ID != "" && call.Function.Name != "" {
					toolNames[call.ID] = call.Function.Name
				}
			}
			out = append(out, raw)
			continue
		}
		if msg.Role != "tool" {
			out = append(out, raw)
			continue
		}

		text := toolResultText(msg.Content)
		if name := toolNames[msg.ToolCallID]; name != "" {
			text = "[" + name + " output]\n" + text
		}
		folded, err := json.Marshal(struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{Role: "user", Content: text})
		if err != nil {
			out = append(out, raw)
			continue
		}
		out = append(out, folded)
		changed = true
	}
	if !changed {
		return body, false
	}

	encodedMessages, err := json.Marshal(out)
	if err != nil {
		return body, false
	}
	doc["messages"] = encodedMessages
	encoded, err := json.Marshal(doc)
	if err != nil {
		return body, false
	}
	return encoded, true
}

// toolResultText flattens a tool message's content to text. The OpenAI wire
// allows a plain string or an array of content parts; both are handled, and an
// empty result keeps an explicit marker rather than becoming silence.
func toolResultText(content json.RawMessage) string {
	if len(content) == 0 {
		return "(no output)"
	}
	var asString string
	if err := json.Unmarshal(content, &asString); err == nil {
		if asString == "" {
			return "(no output)"
		}
		return asString
	}
	var parts []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if err := json.Unmarshal(content, &parts); err == nil {
		var b strings.Builder
		for _, part := range parts {
			if part.Type == "text" {
				b.WriteString(part.Text)
			}
		}
		if b.Len() == 0 {
			return "(no output)"
		}
		return b.String()
	}
	return "(no output)"
}

// itoa avoids pulling strconv in for one call site in this file.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
