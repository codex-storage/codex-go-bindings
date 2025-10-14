package codex

import (
	"bytes"
	"log"
	"os"
	"testing"
)

const expectedCID = "zDvZRwzmAkhzDRPH5EW242gJBNZ2T7aoH2v1fVH66FxXL4kSbvyM"

func TestUploadReader(t *testing.T) {
	codex := newCodexNode(t)
	totalBytes := 0
	finalPercent := 0.0

	buf := bytes.NewBuffer([]byte("Hello World!"))
	len := buf.Len()
	cid, err := codex.UploadReader(UploadOptions{Filepath: "hello.txt", OnProgress: func(read, total int, percent float64, err error) {
		if err != nil {
			log.Fatalf("Error happened during upload: %v\n", err)
		}

		totalBytes = total
		finalPercent = percent
	}}, buf)

	if err != nil {
		t.Fatalf("UploadReader failed: %v", err)
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback read %d bytes but expected %d", totalBytes, len)
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}
}

func TestUploadFile(t *testing.T) {
	codex := newCodexNode(t)
	totalBytes := 0
	finalPercent := 0.0

	stat, err := os.Stat("./testdata/hello.txt")
	if err != nil {
		log.Fatalf("Error happened during file stat: %v\n", err)
	}

	options := UploadOptions{Filepath: "./testdata/hello.txt", OnProgress: func(read, total int, percent float64, err error) {
		if err != nil {
			log.Fatalf("Error happened during upload: %v\n", err)
		}

		totalBytes = total
		finalPercent = percent
	}}

	cid, err := codex.UploadFile(options)
	if err != nil {
		t.Fatalf("UploadReader failed: %v", err)
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}

	if totalBytes != int(stat.Size()) {
		t.Fatalf("UploadReader progress callback read %d bytes but expected %d", totalBytes, int(stat.Size()))
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}
}

func TestUploadFileNoProgress(t *testing.T) {
	codex := newCodexNode(t)

	options := UploadOptions{Filepath: "./testdata/doesnt_exist.txt"}

	cid, err := codex.UploadFile(options)
	if err == nil {
		t.Fatalf("UploadReader should have failed")
	}

	if cid != "" {
		t.Fatalf("Cid should be empty but got %s", cid)
	}
}

func TestManualUpload(t *testing.T) {
	codex := newCodexNode(t)

	sessionId, err := codex.UploadInit(&UploadOptions{Filepath: "hello.txt"})
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	err = codex.UploadChunk(sessionId, []byte("Hello "))
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	err = codex.UploadChunk(sessionId, []byte("World!"))
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	cid, err := codex.UploadFinalize(sessionId)
	if err != nil {
		log.Fatal("Error happened:", err.Error())
	}

	if cid != expectedCID {
		t.Fatalf("UploadReader returned %s but expected %s", cid, expectedCID)
	}
}
