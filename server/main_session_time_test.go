package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestMainSegModTime covers the fix for "the main agent says there are newer
// chats but I don't see them".
//
// Segment 0 is implicit: it has no main_sessions row, so ListMainSessions hands
// back a zero time and the UI rendered "0001-01-01" as its last-activity stamp.
// mainSegModTime fills it from the transcript's mtime instead.
//
// The naming has to match agent/mainagent.go exactly — segment 0 keeps the legacy
// "exp<id>-main" name, later segments are "exp<id>-main-s<seq>". A mismatch here
// would silently fall back to the zero time again, so this pins both shapes.
func TestMainSegModTime(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "transcripts"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	want := time.Date(2026, 10, 26, 2, 10, 33, 0, time.UTC)
	files := map[string]string{
		// segment 0 uses the legacy name, with no "-s0" suffix
		"exp3-main.jsonl":    "seg0",
		"exp3-main-s1.jsonl": "seg1",
		"exp3-main-s2.jsonl": "seg2",
	}
	for name := range files {
		p := filepath.Join(dir, "transcripts", name)
		if err := os.WriteFile(p, []byte("{}\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		if err := os.Chtimes(p, want, want); err != nil {
			t.Fatalf("chtimes %s: %v", name, err)
		}
	}

	s := &Server{m: &Manager{dir: dir}}
	task := &Task{ExpID: 3}

	for _, seg := range []int{0, 1, 2} {
		got, ok := s.mainSegModTime(task, seg)
		if !ok {
			t.Errorf("segment %d: not found", seg)
			continue
		}
		if !got.Equal(want) {
			t.Errorf("segment %d: got %v, want %v", seg, got, want)
		}
	}
}

// TestMainSegModTimeMissingFile makes sure a missing transcript is reported as
// absent rather than inventing a timestamp. The caller then leaves the zero time
// in place, which the frontend already knows how to hide.
func TestMainSegModTimeMissingFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "transcripts"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	s := &Server{m: &Manager{dir: dir}}
	if _, ok := s.mainSegModTime(&Task{ExpID: 99}, 0); ok {
		t.Error("reported a timestamp for a transcript that does not exist")
	}
}

// TestMainSegModTimeDoesNotMatchWrongTask guards the path construction: a task
// whose transcripts exist must not pick up another task's files.
func TestMainSegModTimeDoesNotMatchWrongTask(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "transcripts"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	p := filepath.Join(dir, "transcripts", "exp3-main.jsonl")
	if err := os.WriteFile(p, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	s := &Server{m: &Manager{dir: dir}}
	// exp 4 has no transcript; it must not resolve exp 3's file.
	if _, ok := s.mainSegModTime(&Task{ExpID: 4}, 0); ok {
		t.Error("task 4 picked up task 3's transcript")
	}
}
