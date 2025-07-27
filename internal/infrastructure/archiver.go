package infrastructure

import (
	"23-7-2025/internal/business/interfaces"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type ZipArchiver struct {
}

func NewZipArchiver() *ZipArchiver {
	return &ZipArchiver{}
}

func (za *ZipArchiver) CreateArchive(
	archiveName string, dir string,
) (interfaces.Archive, error) {
	return NewZipArchive(archiveName, dir)
}

func (za *ZipArchiver) Extension() string {
	return ".zip"
}

type ZipArchive struct {
	file      *os.File
	zipWriter *zip.Writer
	mu        sync.Mutex
}

func NewZipArchive(archiveName string, dir string) (*ZipArchive, error) {
	file, err := os.Create(dir + "/" + archiveName + ".zip")
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}

	return &ZipArchive{
		file:      file,
		zipWriter: zip.NewWriter(file),
	}, nil
}

func (za *ZipArchive) AddFile(filename string) error {
	za.mu.Lock()
	defer za.mu.Unlock()

	w, err := za.zipWriter.Create(filepath.Base(filename))
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}

func (za *ZipArchive) Path() string {
	return za.file.Name()
}

func (za *ZipArchive) Close() error {
	err1 := za.zipWriter.Close()
	err2 := za.file.Close()
	return errors.Join(err1, err2)
}
