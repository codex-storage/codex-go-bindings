package codex

import (
	"os"
	"strings"
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

func TestUpdateLogLevel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "codex-log-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	node, err := CodexNew(CodexConfig{
		LogFile: tmpFile.Name(),
	})
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}

	t.Cleanup(func() {
		if err := node.Stop(); err != nil {
			t.Logf("cleanup codex: %v", err)
		}

		if err := node.Destroy(); err != nil {
			t.Logf("cleanup codex: %v", err)
		}
	})

	if err := node.Start(); err != nil {
		t.Fatalf("Failed to start Codex node: %v", err)
	}

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "Started codex node") {
		t.Errorf("Log file does not contain 'Started codex node' %s", string(content))
	}

	if err := node.Stop(); err != nil {
		t.Fatalf("Failed to stop Codex node: %v", err)
	}

	err = node.UpdateLogLevel("ERROR")
	if err != nil {
		t.Fatalf("UpdateLogLevel call failed: %v", err)
	}

	if err := os.WriteFile(tmpFile.Name(), []byte{}, 0644); err != nil {
		t.Fatalf("Failed to clear log file: %v", err)
	}

	err = node.Start()
	if err != nil {
		t.Fatalf("Failed to start Codex node: %v", err)
	}

	content, err = os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if strings.Contains(string(content), "Starting discovery node") {
		t.Errorf("Log file contains 'Starting discovery node'")
	}
}
