package interfaces

type HTTPClient interface {
	DownloadFile(uri string, dir string) (string, error)
}
