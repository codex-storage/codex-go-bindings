package codex

import "testing"

func TestCodexVersion(t *testing.T) {
	node, err := CodexNew(CodexConfig{
		DataDir:   t.TempDir(),
		LogFormat: LogFormatNoColors,
	})
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}
	defer node.Destroy()

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
	node, err := CodexNew(CodexConfig{
		DataDir:   t.TempDir(),
		LogFormat: LogFormatNoColors,
	})
	if err != nil {
		t.Fatalf("Failed to create Codex node: %v", err)
	}
	defer node.Destroy()

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
	node := newCodexNode(t)

	repo, err := node.Repo()
	if err != nil {
		t.Fatalf("Failed to get Codex repo: %v", err)
	}
	if repo == "" {
		t.Fatal("Codex repo is empty")
	}

	t.Logf("Codex repo: %s", repo)
}
