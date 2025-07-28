package services

import (
	"23-7-2025/internal/business/apperrors"
	"23-7-2025/internal/business/dtos"
	"23-7-2025/internal/business/options"
	"23-7-2025/internal/entities"
	"fmt"
	"github.com/google/uuid"
	"path"
	"sync"
)

type TaskService struct {
	options        options.TaskOptions
	archiveService ArchiveServicer

	archivedData map[uuid.UUID]*entities.Task
	data         map[uuid.UUID]*entities.Task
	mu           sync.RWMutex
}

func NewTaskService(
	options options.TaskOptions,
	archiveService ArchiveServicer,
) *TaskService {
	return &TaskService{
		options:        options,
		archiveService: archiveService,
		archivedData:   make(map[uuid.UUID]*entities.Task),
		data:           make(map[uuid.UUID]*entities.Task),
	}
}

func (ts *TaskService) CreateTask() (uuid.UUID, error) {
	task := entities.NewTask()
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if len(ts.data) >= ts.options.MaxNumTasks {
		return uuid.Nil, &apperrors.AppError{Msg: "server is busy"}
	}
	ts.data[task.ID] = task

	return task.ID, nil
}

func (ts *TaskService) List(urlBuilder func(taskID uuid.UUID) string) []dtos.Task {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	tasks := make([]dtos.Task, 0, len(ts.data))
	tasks = toDtos(urlBuilder, tasks, ts.data)
	tasks = toDtos(urlBuilder, tasks, ts.archivedData)

	return tasks
}

func (ts *TaskService) AddResource(
	taskID uuid.UUID, resource entities.Resource,
) error {
	ext := path.Ext(resource.URI)
	if _, ok := ts.options.AllowedFileExtensions[ext]; !ok {
		return &apperrors.AppError{Msg: "file extension not allowed"}
	}

	ts.mu.RLock()
	task, ok := ts.data[taskID]
	ts.mu.RUnlock()

	if !ok {
		return &apperrors.NotFoundError{Msg: "task not found"}
	}

	task.Mu.Lock()
	defer task.Mu.Unlock()
	if len(task.Resources) >= ts.options.MaxNumResources {
		return &apperrors.AppError{
			Msg: fmt.Sprintf(
				"task max number resources reached, %d", ts.options.MaxNumResources,
			),
		}
	}
	task.Resources = append(task.Resources, resource)

	return nil
}

func (ts *TaskService) GetTaskStatus(
	urlBuilder func(taskID uuid.UUID) string,
	taskID uuid.UUID,
) (dtos.Task, error) {
	ts.mu.RLock()
	task, ok := ts.getTask(taskID)
	ts.mu.RUnlock()
	if !ok {
		return dtos.Task{}, &apperrors.NotFoundError{Msg: "task not found"}
	}

	task.Mu.Lock()
	defer task.Mu.Unlock()

	if task.Status == entities.StatusPending {
		if len(task.Resources) == ts.options.MaxNumResources {
			archive, err := ts.archiveService.Archive(task)
			if err != nil {
				return dtos.Task{}, fmt.Errorf("failed to archive task: %w", err)
			}
			task.ArchivePath = archive
			task.Status = entities.StatusArchived
			ts.archivedData[task.ID] = task
			delete(ts.data, task.ID)
		}
	}

	taskDTO := dtos.Task{
		ID:    task.ID.String(),
		State: dtos.NewTaskStatus(urlBuilder, task),
	}
	return taskDTO, nil
}

func (ts *TaskService) getTask(id uuid.UUID) (*entities.Task, bool) {
	ts.mu.RLock()
	task, ok := ts.data[id]
	if !ok {
		task, ok = ts.archivedData[id]
	}
	ts.mu.RUnlock()
	return task, ok
}

func toDtos(
	urlBuilder func(taskID uuid.UUID) string,
	out []dtos.Task, tasks map[uuid.UUID]*entities.Task,
) []dtos.Task {
	for _, task := range tasks {
		task.Mu.RLock()
		task.Mu.RUnlock()

		out = append(
			out, dtos.Task{
				ID:    task.ID.String(),
				State: dtos.NewTaskStatus(urlBuilder, task),
			},
		)
	}

	return out
}
