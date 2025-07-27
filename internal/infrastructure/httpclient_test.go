package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHTTPClient_DownloadFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testContent := "file's content"
	fileName := "testfile.txt"

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/"+fileName {
					w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
					_, err = w.Write([]byte(testContent))
					if err != nil {
						t.Fatal(err)
					}
				} else {
					http.NotFound(w, r)
				}
			},
		),
	)
	defer ts.Close()

	downloadURL := ts.URL + "/" + fileName
	client := NewHTTPClient()

	downloadedFile, err := client.DownloadFile(downloadURL, tempDir)
	if err != nil {
		t.Errorf("download error: %v", err)
	}

	_, err = os.Stat(downloadedFile)
	if os.IsNotExist(err) {
		t.Errorf("downloaded file doesn't exist: %v", err)
	} else if err != nil {
		t.Errorf("error statting downloaded file: %v", err)
	}

	downloadedContent, err := os.ReadFile(downloadedFile)
	if err != nil {
		t.Errorf("error reading downloaded file: %v", err)
	}
	if string(downloadedContent) != testContent {
		t.Errorf("downloaded file content doesn't match: %s", string(downloadedContent))
	}

	t.Run(
		"non-existent file",
		func(t *testing.T) {
			badURL := ts.URL + "/non-existent-file.txt"
			_, err := client.DownloadFile(badURL, tempDir)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		},
	)

	t.Run(
		"invalid dir",
		func(t *testing.T) {
			invalidDir := "invalid-dir/2281337"
			_, err := client.DownloadFile(downloadURL, invalidDir)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		},
	)
}
