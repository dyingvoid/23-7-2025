package services

import (
	"23-7-2025/internal/business/interfaces"
	"23-7-2025/internal/entities"
	"fmt"
	"os"
	"sync"
)

type ArchiveServicer interface {
	ArchivePath(task *entities.Task) string
	Archive(task *entities.Task) (string, error)
}

type ArchiveService struct {
	archiver        interfaces.Archivator
	resourceService *ResourceService
	archiveDir      string
}

func NewArchiveService(
	archiver interfaces.Archivator, resourceService *ResourceService, dir string,
) *ArchiveService {
	return &ArchiveService{
		archiver:        archiver,
		resourceService: resourceService,
		archiveDir:      dir,
	}
}

func (as *ArchiveService) ArchivePath(task *entities.Task) string {
	return as.archiveDir + task.ID.String() + as.archiver.Extension()
}

func (as *ArchiveService) Archive(task *entities.Task) (string, error) {
	tempDir, err := os.MkdirTemp("", "download_"+task.ID.String()+"_*")
	if err != nil {
		return "", fmt.Errorf("couldn't create temp archiveDir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archive, err := as.archiver.CreateArchive(task.ID.String(), as.archiveDir)
	if err != nil {
		return "", fmt.Errorf("couldn't create archive: %w", err)
	}
	defer archive.Close()

	wg := sync.WaitGroup{}
	for i := 0; i < len(task.Resources); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resource := as.resourceService.DownloadResource(task.Resources[i], tempDir)
			if resource.Error != nil || !resource.Downloaded {
				return
			}
			err = archive.AddFile(resource.Filename)
			if err != nil {
				resource.Archived = false
				resource.Error = fmt.Errorf("couldn't add file to archive: %w", err)
			}

			task.Resources[i] = resource
		}()
	}
	wg.Wait()

	return archive.Path(), nil
}
