package interfaces

type HTTPClienter interface {
	DownloadFile(uri string, dir string) (string, error)
}
