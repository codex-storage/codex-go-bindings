package codex

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func uploadHelper(t *testing.T, codex *CodexNode) (string, int) {
	t.Helper()

	buf := bytes.NewBuffer([]byte("Hello World!"))
	len := buf.Len()
	cid, err := codex.UploadReader(UploadOptions{filepath: "hello.txt"}, buf)
	if err != nil {
		t.Fatalf("Error happened during upload: %v\n", err)
	}

	return cid, len
}

func TestDownloadStream(t *testing.T) {
	start := true
	codex := newCodexNode(t, start)
	cid, len := uploadHelper(t, codex)

	f, err := os.Create("testdata/hello.downloaded.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	totalBytes := 0
	finalPercent := 0.0
	opt := DownloadStreamOptions{
		writer:      f,
		datasetSize: len,
		filepath:    "testdata/hello.downloaded.writer.txt",
		onProgress: func(read, total int, percent float64, err error) {
			if err != nil {
				t.Fatalf("Error happening during download: %v\n", err)
			}

			totalBytes = total
			finalPercent = percent
		},
	}

	if err := codex.DownloadStream(cid, opt); err != nil {
		t.Fatal("Error happened:", err.Error())
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback total bytes %d but expected %d", totalBytes, len)
	}

	data, err := os.ReadFile("testdata/hello.writer.txt")
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hello World!" {
		t.Fatalf("Downloaded content does not match, expected Hello World! got %s", data)
	}
}

func TestDownloadStreamWithAutosize(t *testing.T) {
	start := true
	codex := newCodexNode(t, start)
	cid, len := uploadHelper(t, codex)

	totalBytes := 0
	finalPercent := 0.0
	opt := DownloadStreamOptions{
		datasetSizeAuto: true,
		onProgress: func(read, total int, percent float64, err error) {
			if err != nil {
				t.Fatalf("Error happening during download: %v\n", err)
			}

			totalBytes = total
			finalPercent = percent
		},
	}

	if err := codex.DownloadStream(cid, opt); err != nil {
		t.Fatal("Error happened:", err.Error())
	}

	if finalPercent != 100.0 {
		t.Fatalf("UploadReader progress callback final percent %.2f but expected 100.0", finalPercent)
	}

	if totalBytes != len {
		t.Fatalf("UploadReader progress callback total bytes %d but expected %d", totalBytes, len)
	}
}

func TestDownloadManual(t *testing.T) {
	start := true
	codex := newCodexNode(t, start)
	cid, _ := uploadHelper(t, codex)

	if err := codex.DownloadInit(cid, DownloadInitOptions{}); err != nil {
		t.Fatal("Error when initializing download:", err)
	}

	var b strings.Builder
	if chunk, err := codex.DownloadChunk(cid); err != nil {
		t.Fatal("Error when downloading chunk:", err)
	} else {
		b.Write(chunk)
	}

	data := b.String()
	if data != "Hello World!" {
		t.Fatalf("Expected data was \"Hello World!\" got %s", data)
	}

	if err := codex.DownloadCancel(cid); err != nil {
		t.Fatalf("Error when cancelling the download %s", err)
	}
}
