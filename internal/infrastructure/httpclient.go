package infrastructure

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: http.DefaultClient,
	}
}

func (c *HTTPClient) DownloadFile(uri string, dir string) (string, error) {
	resp, err := c.client.Get(uri)
	if err != nil {
		return "", fmt.Errorf("error while getting uri: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	outFile, err := os.Create(dir + "/" + path.Base(uri))
	if err != nil {
		return "", fmt.Errorf("error while creating file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while copying file: %w", err)
	}

	return outFile.Name(), nil
}
