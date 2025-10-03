package codex

import (
	"testing"
)

func TestDebug(t *testing.T) {
	codex := newCodexNode(t)

	info, err := codex.Debug()
	if err != nil {
		t.Fatalf("Debug call failed: %v", err)
	}
	if info.ID == "" {
		t.Error("Debug info ID is empty")
	}
	if info.Spr == "" {
		t.Error("Debug info Spr is empty")
	}
	if len(info.AnnounceAddresses) == 0 {
		t.Error("Debug info AnnounceAddresses is empty")
	}
}
