package tests

import (
	"23-7-2025/internal/business/services"
	"23-7-2025/internal/entities"
	"23-7-2025/internal/infrastructure"
	"archive/zip"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

func TestArchiving(t *testing.T) {
	tempDir := t.TempDir()
	zipArchiver := infrastructure.NewZipArchiver()
	client := infrastructure.NewHTTPClient()
	resourceService := services.NewResourceService(client, tempDir)
	archiveService := services.NewArchiveService(zipArchiver, resourceService, tempDir)

	t.Run(
		"diff names success",
		func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fileName := strings.TrimPrefix(r.URL.Path, "/")
						if fileName == "" {
							fileName = "default.txt"
						}
						w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
						_, err := w.Write([]byte(fileName))
						if err != nil {
							t.Fatal(err)
						}
					},
				),
			)
			defer ts.Close()

			task := entities.NewTask()
			task.Resources = []entities.Resource{
				entities.NewResource(ts.URL + "/test1.txt"),
				entities.NewResource(ts.URL + "/test2.txt"),
				entities.NewResource(ts.URL + "/test3.txt"),
			}
			p, err := archiveService.Archive(task)
			if err != nil {
				t.Fatal(err)
			}

			zor, err := zip.OpenReader(p)
			if err != nil {
				t.Fatal(err)
			}
			defer zor.Close()

			expectedFiles := map[string]bool{}
			for _, r := range task.Resources {
				expectedFiles[path.Base(r.URI)] = false
			}

			for _, f := range zor.File {
				if _, ok := expectedFiles[f.Name]; ok {
					expectedFiles[f.Name] = true

					rc, err := f.Open()
					assert.NoErrorf(t, err, "failed to open file %s in archive", f.Name)
					if err == nil {
						rc.Close()
					}
				}
			}

			for fileName, found := range expectedFiles {
				assert.Truef(t, found, "file %s not found in archive", fileName)
			}
		},
	)
	t.Run(
		"same names success",
		func(t *testing.T) {
			ts := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						fileName := strings.TrimPrefix(r.URL.Path, "/")
						if fileName == "" {
							fileName = "default.txt"
						}
						w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
						_, err := w.Write([]byte(fileName))
						if err != nil {
							t.Fatal(err)
						}
					},
				),
			)
			defer ts.Close()

			task := entities.NewTask()
			task.Resources = []entities.Resource{
				entities.NewResource(ts.URL + "/test1.txt"),
				entities.NewResource(ts.URL + "/test1.txt"),
				entities.NewResource(ts.URL + "/test1.txt"),
			}
			p, err := archiveService.Archive(task)
			if err != nil {
				t.Fatal(err)
			}

			zor, err := zip.OpenReader(p)
			if err != nil {
				t.Fatal(err)
			}
			defer zor.Close()

			expectedFiles := map[string]bool{}
			for _, r := range task.Resources {
				expectedFiles[path.Base(r.URI)] = false
			}

			assert.True(t, len(zor.File) == len(task.Resources))
			for _, f := range zor.File {
				if _, ok := expectedFiles[f.Name]; ok {
					expectedFiles[f.Name] = true

					rc, err := f.Open()
					assert.NoErrorf(t, err, "failed to open file %s in archive", f.Name)
					if err == nil {
						rc.Close()
					}
				}
			}

			for fileName, found := range expectedFiles {
				assert.Truef(t, found, "file %s not found in archive", fileName)
			}
		},
	)
}
