package interfaces

type Interfaces struct {
	Archiver   Archivator
	HTTPClient HTTPClient
}

func New(archiver Archivator, httpClient HTTPClient) *Interfaces {
	return &Interfaces{
		Archiver:   archiver,
		HTTPClient: httpClient,
	}
}
