package codex

import (
	"os"
	"strings"
	"testing"
	"time"
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

	node, err := New(Config{
		LogFile:        tmpFile.Name(),
		MetricsEnabled: false,
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

func TestCodexPeerDebug(t *testing.T) {
	var bootstrap, node1, node2 *CodexNode
	var err error

	t.Cleanup(func() {
		if bootstrap != nil {
			if err := bootstrap.Stop(); err != nil {
				t.Logf("cleanup bootstrap: %v", err)
			}

			if err := bootstrap.Destroy(); err != nil {
				t.Logf("cleanup bootstrap: %v", err)
			}
		}
		if node1 != nil {
			if err := node1.Stop(); err != nil {
				t.Logf("cleanup node1: %v", err)
			}

			if err := node1.Destroy(); err != nil {
				t.Logf("cleanup node1: %v", err)
			}
		}
		if node2 != nil {
			if err := node2.Stop(); err != nil {
				t.Logf("cleanup node2: %v", err)
			}

			if err := node2.Destroy(); err != nil {
				t.Logf("cleanup node2: %v", err)
			}
		}
	})

	bootstrap, err = New(Config{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		DiscoveryPort:  8092,
	})
	if err != nil {
		t.Fatalf("Failed to create bootstrap: %v", err)
	}

	if err := bootstrap.Start(); err != nil {
		t.Fatalf("Failed to start bootstrap: %v", err)
	}

	spr, err := bootstrap.Spr()
	if err != nil {
		t.Fatalf("Failed to get bootstrap spr: %v", err)
	}

	bootstrapNodes := []string{spr}

	node1, err = New(Config{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		DiscoveryPort:  8090,
		BootstrapNodes: bootstrapNodes,
	})
	if err != nil {
		t.Fatalf("Failed to create codex: %v", err)
	}

	if err := node1.Start(); err != nil {
		t.Fatalf("Failed to start codex: %v", err)
	}

	node2, err = New(Config{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		DiscoveryPort:  8091,
		BootstrapNodes: bootstrapNodes,
	})
	if err != nil {
		t.Fatalf("Failed to create codex2: %v", err)
	}

	if err := node2.Start(); err != nil {
		t.Fatalf("Failed to start codex2: %v", err)
	}

	peerId, err := node2.PeerId()
	if err != nil {
		t.Fatal(err)
	}

	var record PeerRecord
	for range 10 {
		record, err = node1.CodexPeerDebug(peerId)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if record.PeerId == "" {
		t.Fatalf("CodexPeerDebug call failed: %v", err)
	}
	if record.PeerId == "" {
		t.Error("CodexPeerDebug info PeerId is empty")
	}
	if record.SeqNo == 0 {
		t.Error("CodexPeerDebug info SeqNo is empty")
	}
	if len(record.Addresses) == 0 {
		t.Error("CodexPeerDebug info Addresses is empty")
	}
	if record.PeerId != peerId {
		t.Errorf("CodexPeerDebug info PeerId (%s) does not match requested PeerId (%s)", record.PeerId, peerId)
	}
}
