package interfaces

type Interfaces struct {
	Archiver   Archiver
	HTTPClient HTTPClienter
}

func New(archiver Archiver, httpClient HTTPClienter) *Interfaces {
	return &Interfaces{
		Archiver:   archiver,
		HTTPClient: httpClient,
	}
}
