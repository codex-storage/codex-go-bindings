package codex

import (
	"log"
	"testing"
)

func TestConnectWithAddress(t *testing.T) {
	var node1, node2 *CodexNode
	var err error

	t.Cleanup(func() {
		if node1 != nil {
			if err := node1.Stop(); err != nil {
				t.Logf("cleanup codex1: %v", err)
			}

			if err := node1.Destroy(); err != nil {
				t.Logf("cleanup codex1: %v", err)
			}
		}

		if node2 != nil {
			if err := node2.Stop(); err != nil {
				t.Logf("cleanup codex2: %v", err)
			}

			if err := node2.Destroy(); err != nil {
				t.Logf("cleanup codex2: %v", err)
			}
		}
	})

	node1, err = CodexNew(CodexConfig{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		DiscoveryPort:  8090,
	})
	if err != nil {
		t.Fatalf("Failed to create codex1: %v", err)
	}

	if err := node1.Start(); err != nil {
		t.Fatalf("Failed to start codex1: %v", err)
	}

	node2, err = CodexNew(CodexConfig{
		DataDir:        t.TempDir(),
		LogFormat:      LogFormatNoColors,
		MetricsEnabled: false,
		DiscoveryPort:  8091,
	})
	if err != nil {
		t.Fatalf("Failed to create codex2: %v", err)
	}

	if err := node2.Start(); err != nil {
		t.Fatalf("Failed to start codex2: %v", err)
	}

	info2, err := node2.Debug()
	if err != nil {
		t.Fatal(err)
	}

	if err := node1.Connect(info2.ID, info2.Addrs); err != nil {
		t.Fatalf("connect failed: %v", err)
	}
}

func TestCodexWithPeerId(t *testing.T) {
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

	bootstrap, err = CodexNew(CodexConfig{
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

	node1, err = CodexNew(CodexConfig{
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

	node2, err = CodexNew(CodexConfig{
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

	if err := node1.Connect(peerId, []string{}); err != nil {
		log.Println(err)
	}
}
