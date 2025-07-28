package services

import (
	"23-7-2025/internal/business/apperrors"
	"23-7-2025/internal/business/options"
	"23-7-2025/internal/entities"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateTask_Success(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	ts := NewTaskService(opts, mockArchive)

	id, err := ts.CreateTask()
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreateTask_TooManyTasks(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	opts.MaxNumTasks = 1
	ts := NewTaskService(opts, mockArchive)

	_, err1 := ts.CreateTask()
	assert.NoError(t, err1)
	_, err2 := ts.CreateTask()
	assert.Error(t, err2)
	assert.IsType(t, &apperrors.AppError{}, err2)
}

func TestAddResource_ValidResource(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	ts := NewTaskService(makeOptions(), mockArchive)
	id, _ := ts.CreateTask()
	resource := makeResource("file.jpeg")

	err := ts.AddResource(id, resource)
	assert.NoError(t, err)

	task := ts.data[id]
	assert.Len(t, task.Resources, 1)
	assert.Equal(t, resource, task.Resources[0])
}

func TestAddResource_DisallowedExtension(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	ts := NewTaskService(makeOptions(), mockArchive)
	id, _ := ts.CreateTask()
	resource := makeResource("file.exe")

	err := ts.AddResource(id, resource)
	assert.Error(t, err)
	assert.IsType(t, &apperrors.AppError{}, err)
}

func TestAddResource_UnknownTask(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	ts := NewTaskService(makeOptions(), mockArchive)

	resource := makeResource("file.pdf")
	err := ts.AddResource(uuid.New(), resource)
	assert.Error(t, err)
	assert.IsType(t, &apperrors.NotFoundError{}, err)
}

func TestAddResource_MaxResources(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	ts := NewTaskService(opts, mockArchive)
	id, _ := ts.CreateTask()

	_ = ts.AddResource(id, makeResource("a.jpg"))
	_ = ts.AddResource(id, makeResource("b.jpg"))
	err := ts.AddResource(id, makeResource("c.jpg"))
	assert.Error(t, err)
	assert.IsType(t, &apperrors.AppError{}, err)
}

func TestGetTaskStatus_Archive(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	ts := NewTaskService(opts, mockArchive)
	id, _ := ts.CreateTask()

	// Add resources to reach max cap
	err := ts.AddResource(id, makeResource("a.jpeg"))
	assert.NoError(t, err)
	err = ts.AddResource(id, makeResource("b.pdf"))
	assert.NoError(t, err)

	mockArchive.On("Archive", mock.Anything).Return("archive/path", nil)

	task, err := ts.GetTaskStatus(urlBuilder, id)
	assert.NoError(t, err)
	assert.Equal(t, entities.StatusArchived.String(), task.State.Status)
	mockArchive.AssertCalled(t, "Archive", mock.Anything)
}

func TestGetTaskStatus_NotFound(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	ts := NewTaskService(makeOptions(), mockArchive)

	_, err := ts.GetTaskStatus(urlBuilder, uuid.New())
	assert.Error(t, err)
	assert.IsType(t, &apperrors.NotFoundError{}, err)
}

func TestGetTaskStatus_ArchiveFail(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	ts := NewTaskService(opts, mockArchive)
	id, _ := ts.CreateTask()

	_ = ts.AddResource(id, makeResource("a.jpeg"))
	_ = ts.AddResource(id, makeResource("b.jpeg"))
	mockArchive.On("Archive", mock.Anything).
		Return("", errors.New("fail"))

	_, err := ts.GetTaskStatus(urlBuilder, id)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to archive task")
}

func TestList_ReturnsActiveAndArchivedTasks(t *testing.T) {
	mockArchive := new(MockArchiveServicer)
	opts := makeOptions()
	ts := NewTaskService(opts, mockArchive)
	id, err := ts.CreateTask()
	assert.NoError(t, err)
	_, err = ts.CreateTask()
	assert.NoError(t, err)

	err = ts.AddResource(id, makeResource("a.jpeg"))
	assert.NoError(t, err)
	err = ts.AddResource(id, makeResource("b.pdf"))
	assert.NoError(t, err)
	mockArchive.On("Archive", mock.Anything).
		Return("archive/path", nil)
	_, _ = ts.GetTaskStatus(urlBuilder, id) // Moves to archived

	tasks := ts.List(urlBuilder)

	assert.Len(t, tasks, 2)
	hasArchived := false
	for _, task := range tasks {
		if task.State.Status == entities.StatusArchived.String() {
			hasArchived = true
		}
	}
	assert.True(t, hasArchived)
}

type MockArchiveServicer struct {
	mock.Mock
}

func (m *MockArchiveServicer) ArchivePath(taskID uuid.UUID) string {
	args := m.Called(taskID)
	return args.String(0)
}

func (m *MockArchiveServicer) Archive(task *entities.Task) (string, error) {
	args := m.Called(task)
	return args.String(0), args.Error(1)
}

func makeOptions() options.TaskOptions {
	return options.TaskOptions{
		MaxNumTasks:           2,
		MaxNumResources:       2,
		AllowedFileExtensions: map[string]struct{}{".jpeg": {}, ".pdf": {}},
	}
}

func urlBuilder(taskID uuid.UUID) string {
	return ""
}

func makeResource(uri string) entities.Resource {
	return entities.Resource{URI: uri}
}
