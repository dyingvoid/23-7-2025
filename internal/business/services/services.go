package services

import (
	"23-7-2025/internal/business/interfaces"
	"23-7-2025/internal/business/options"
)

type Services struct {
	Task           *TaskService
	ArchiveService *ArchiveService
}

func New(i *interfaces.Interfaces, opts options.TaskOptions) *Services {
	resourceService := NewResourceService(i.HTTPClient)
	archiveService := NewArchiveService(i.Archiver, resourceService, opts.FileDir)
	return &Services{
		Task:           NewTaskService(opts, archiveService),
		ArchiveService: archiveService,
	}
}
