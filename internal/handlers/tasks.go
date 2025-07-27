package handlers

import (
	"23-7-2025/internal/business/dtos"
	"23-7-2025/internal/business/services"
	"23-7-2025/internal/entities"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TasksHandler struct {
	taskService *services.TaskService
}

type ListResponse struct {
	Tasks []dtos.Task `json:"tasks"`
}

// ListTasks
// @Summary List all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {object} ListResponse
// @Router /api/v1/tasks [get]
func (th *TasksHandler) ListTasks(c echo.Context) error {
	tasks := th.taskService.List()
	return c.JSON(
		http.StatusOK, ListResponse{
			Tasks: tasks,
		},
	)
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

// CreateTask
// @Summary Create a new task
// @Description Creates a new task and returns the task ID
// @Tags tasks
// @Produce json
// @Success 200 {object} CreateTaskResponse
// @Router /api/v1/tasks [post]
func (th *TasksHandler) CreateTask(c echo.Context) error {
	id, err := th.taskService.CreateTask()
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK, CreateTaskResponse{
			ID: id.String(),
		},
	)
}

type GetTaskRequest struct {
	ID string `param:"id"`
}
type GetTaskResponse struct {
	Task dtos.Task `json:"task"`
}

// GetTask
// @Summary Gets task state
// @Description Gets task state by ID, if number resources == X, archives resources and returns archive link
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/tasks/{id} [get]
func (th *TasksHandler) GetTask(c echo.Context) error {
	req := new(GetTaskRequest)
	if err := c.Bind(req); err != nil {
		return badRequest(c, err)
	}
	taskID, err := uuid.Parse(req.ID)
	if err != nil {
		return badRequest(c, err)
	}
	task, err := th.taskService.GetTaskStatus(taskID)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK, GetTaskResponse{
			Task: task,
		},
	)
}

type AddResourceRequest struct {
	TaskID      string `param:"id" swaggerignore:"true"`
	ResourceURI string `json:"resource_uri"`
}

// AddResource
// @Summary Adds resource to a task
// @Description Adds resource to a task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Param resource body AddResourceRequest true "ResourceURI object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/tasks/{id}/resources [post]
func (th *TasksHandler) AddResource(c echo.Context) error {
	req := new(AddResourceRequest)
	if err := c.Bind(req); err != nil {
		return badRequest(c, err)
	}
	taskID, err := uuid.Parse(req.TaskID)
	if err != nil {
		return badRequest(c, err)
	}

	err = th.taskService.AddResource(
		taskID, entities.NewResource(req.ResourceURI),
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, map[string]string{})
}

// TODO to utils
func badRequest(c echo.Context, err error) error {
	return c.JSON(
		http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		},
	)
}
