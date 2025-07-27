package infrastructure

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZipArchive(t *testing.T) {
	tmpDir := t.TempDir()
	testData := []byte("hello, zip!")
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatalf("could not write test data: %v", err)
	}

	zipArchive, err := NewZipArchive("test_archive", tmpDir)
	if err != nil {
		t.Fatalf("failed to create zip archive: %v", err)
	}

	if err := zipArchive.AddFile(testFile); err != nil {
		t.Fatalf("failed to add file: %v", err)
	}

	if err := zipArchive.Close(); err != nil {
		t.Fatalf("closing archive failed: %v", err)
	}

	archivePath := filepath.Join(tmpDir, "test_archive.zip")
	zor, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("could not open archive: %v", err)
	}
	defer zor.Close()

	var found bool
	for _, f := range zor.File {
		if f.Name == "test.txt" {
			found = true
			r, err := f.Open()
			if err != nil {
				t.Fatalf("could not open file in archive: %v", err)
			}
			defer r.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			if !strings.EqualFold(string(buf.Bytes()), string(testData)) {
				t.Fatalf("file content mismatch: got %s, want %s", buf.Bytes(), testData)
			}
		}
	}
	if !found {
		t.Fatalf("file 'test.txt' not found in archive")
	}
}
