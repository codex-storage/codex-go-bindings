package codex

import (
	"testing"
)

func newCodexNode(t *testing.T) *CodexNode {
	node, err := CodexNew(CodexConfig{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
	})
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}

	err = node.Start()
	if err != nil {
		t.Fatalf("Failed to start Codex node: %v", err)
	}

	t.Cleanup(func() {
		if err := node.Stop(); err != nil {
			t.Logf("cleanup codex: %v", err)
		}

		if err := node.Destroy(); err != nil {
			t.Logf("cleanup codex: %v", err)
		}
	})

	return node
}
