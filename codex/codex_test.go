package codex

import "testing"

func TestCodexVersion(t *testing.T) {
	start := false
	node := newCodexNode(t, start)

	version, err := node.Version()
	if err != nil {
		t.Fatalf("Failed to get Codex version: %v", err)
	}
	if version == "" {
		t.Fatal("Codex version is empty")
	}

	t.Logf("Codex version: %s", version)
}

func TestCodexRevision(t *testing.T) {
	start := false
	node := newCodexNode(t, start)

	revision, err := node.Revision()
	if err != nil {
		t.Fatalf("Failed to get Codex revision: %v", err)
	}
	if revision == "" {
		t.Fatal("Codex revision is empty")
	}

	t.Logf("Codex revision: %s", revision)
}

func TestCodexRepo(t *testing.T) {
	start := true
	node := newCodexNode(t, start)

	repo, err := node.Repo()
	if err != nil {
		t.Fatalf("Failed to get Codex repo: %v", err)
	}
	if repo == "" {
		t.Fatal("Codex repo is empty")
	}

	t.Logf("Codex repo: %s", repo)
}

func TestSpr(t *testing.T) {
	start := true
	node := newCodexNode(t, start)

	spr, err := node.Spr()
	if err != nil {
		t.Fatalf("Failed to get Codex SPR: %v", err)
	}
	if spr == "" {
		t.Fatal("Codex SPR is empty")
	}

	t.Logf("Codex SPR: %s", spr)
}

func TestPeerId(t *testing.T) {
	start := true
	node := newCodexNode(t, start)

	peerId, err := node.PeerId()
	if err != nil {
		t.Fatalf("Failed to get Codex PeerId: %v", err)
	}
	if peerId == "" {
		t.Fatal("Codex PeerId is empty")
	}

	t.Logf("Codex PeerId: %s", peerId)
}
