package handlers

import (
	"23-7-2025/internal/business/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
)

type Handlers struct {
	Tasks *TasksHandler
}

func New(s *services.Services) *Handlers {
	return &Handlers{
		Tasks: &TasksHandler{taskService: s.Task},
	}
}

func SetDefault(e *echo.Echo) {
	e.Static("/", "static")
	e.GET("/healthcheck", HealthCheckHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}

func SetAPI(e *echo.Echo, h *Handlers) {
	g := e.Group("/api/v1")

	g.GET("/tasks", h.Tasks.ListTasks)
	g.GET("/tasks/:id", h.Tasks.GetTask)
	g.POST("/tasks", h.Tasks.CreateTask)

	g.POST("/tasks/:id/resources", h.Tasks.AddResource)
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
