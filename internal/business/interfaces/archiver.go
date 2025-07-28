package interfaces

type Archiver interface {
	CreateArchive(archiveName string, dir string) (Archive, error)
	Extension() string
}

type Archive interface {
	AddFile(filename string) error
	Path() string
	Close() error
}
