package handlers

import (
	"23-7-2025/internal/business/apperrors"
	"23-7-2025/internal/business/services"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"net/http"
)

type Handlers struct {
	Tasks *TasksHandler
}

func New(s *services.Services) *Handlers {
	return &Handlers{
		Tasks: NewTasksHandler(s.Task, s.ArchiveService),
	}
}

func SetDefault(e *echo.Echo) {
	e.Static("/", "static")
	e.GET("/healthcheck", HealthCheckHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HTTPErrorHandler = errorHandler()
}

func errorHandler() func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		c.Logger().Error(err)
		if c.Response().Committed {
			return
		}
		var httpErr *echo.HTTPError
		if errors.As(err, &httpErr) {
			return
		}

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			sendError(c, http.StatusBadRequest, appErr.Msg)
			return
		}
		var notFoundErr *apperrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			sendError(c, http.StatusNotFound, notFoundErr.Msg)
			return
		}
		var busyErr *apperrors.ServerBusyError
		if errors.As(err, &busyErr) {
			sendError(c, http.StatusServiceUnavailable, "servery busy")
			return
		}

		sendError(c, http.StatusInternalServerError, "Internal server error")
	}
}

func SetAPI(e *echo.Echo, h *Handlers) {
	g := e.Group("/api/v1")

	g.GET("/tasks", h.Tasks.ListTasks)
	g.GET("/tasks/:id", h.Tasks.GetTask)
	g.POST("/tasks", h.Tasks.CreateTask)

	g.POST("/tasks/:id/resources", h.Tasks.AddResource)

	g.GET("/tasks/:id/archive", h.Tasks.GetArchive)
}

func Echo() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{"*"},
			},
		),
	)

	return e
}

func sendError(c echo.Context, code int, msg string) {
	c.JSON(code, map[string]string{"error": msg})
}
