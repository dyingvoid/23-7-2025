package services

import (
	"23-7-2025/internal/business/interfaces"
	"23-7-2025/internal/entities"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

type (
	ArchiveServicer interface {
		ArchivePath(taskID uuid.UUID) string
		Archive(task *entities.Task) (string, error)
	}

	ArchiveService struct {
		archiver        interfaces.Archiver
		resourceService *ResourceService
		archiveDir      string
	}
)

func NewArchiveService(
	archiver interfaces.Archiver, resourceService *ResourceService, dir string,
) *ArchiveService {
	return &ArchiveService{
		archiver:        archiver,
		resourceService: resourceService,
		archiveDir:      dir,
	}
}

func (as *ArchiveService) ArchivePath(taskID uuid.UUID) string {
	return as.archiveDir + "/" + taskID.String() + as.archiver.Extension()
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
				task.Resources[i] = resource
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
