package codex

import "testing"

func TestManifests(t *testing.T) {
	codex := newCodexNode(t)

	manifests, err := codex.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) != 0 {
		t.Fatal("expected manifests to be empty")
	}

	cid, _ := uploadHelper(t, codex)

	manifests, err = codex.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) == 0 {
		t.Fatal("expected manifests to be non-empty")
	}

	for _, m := range manifests {
		if m.Cid != cid {
			t.Errorf("expected cid %q, got %q", cid, m.Cid)
		}
	}
}

func TestSpace(t *testing.T) {
	codex := newCodexNode(t)

	space, err := codex.Space()
	if err != nil {
		t.Fatal(err)
	}

	if space.TotalBlocks != 0 {
		t.Fatal("expected total blocks to be non-zero")
	}

	if space.QuotaMaxBytes == 0 {
		t.Fatal("expected quota max bytes to be non-zero")
	}

	if space.QuotaUsedBytes != 0 {
		t.Fatal("expected quota used bytes to be non-zero")
	}

	if space.QuotaReservedBytes != 0 {
		t.Fatal("expected quota reserved bytes to be non-zero")
	}

	uploadHelper(t, codex)

	space, err = codex.Space()
	if err != nil {
		t.Fatal(err)
	}

	if space.TotalBlocks == 0 {
		t.Fatal("expected total blocks to be non-zero after upload")
	}

	if space.QuotaUsedBytes == 0 {
		t.Fatal("expected quota used bytes to be non-zero after upload")
	}
}

func TestFetch(t *testing.T) {
	codex := newCodexNode(t)

	cid, _ := uploadHelper(t, codex)

	_, err := codex.Fetch(cid)
	if err != nil {
		t.Fatal("expected error when fetching non-existent manifest")
	}
}

func TestFetchCidDoesNotExist(t *testing.T) {
	codex := newCodexNode(t, Config{BlockRetries: 1})

	_, err := codex.Fetch("bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku")
	if err == nil {
		t.Fatal("expected error when fetching non-existent manifest")
	}
}

func TestDelete(t *testing.T) {
	codex := newCodexNode(t)

	cid, _ := uploadHelper(t, codex)

	manifests, err := codex.Manifests()
	if err != nil {
		t.Fatal(err)
	}
	if len(manifests) != 1 {
		t.Fatal("expected manifests to be empty after deletion")
	}

	err = codex.Delete(cid)
	if err != nil {
		t.Fatal(err)
	}

	manifests, err = codex.Manifests()
	if err != nil {
		t.Fatal(err)
	}

	if len(manifests) != 0 {
		t.Fatal("expected manifests to be empty after deletion")
	}
}

func TestExists(t *testing.T) {
	codex := newCodexNode(t)

	cid, _ := uploadHelper(t, codex)

	exists, err := codex.Exists(cid)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected cid to exist")
	}

	err = codex.Delete(cid)
	if err != nil {
		t.Fatal(err)
	}

	exists, err = codex.Exists(cid)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("expected cid to not exist after deletion")
	}
}
