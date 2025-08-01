package main

import (
	"23-7-2025/config"
	_ "23-7-2025/docs"
	"23-7-2025/internal/business/interfaces"
	"23-7-2025/internal/business/options"
	"23-7-2025/internal/business/services"
	"23-7-2025/internal/handlers"
	"23-7-2025/internal/infrastructure"
	"errors"
	"os"
	"strconv"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
func main() {
	// TODO to env
	cfg := config.RequireConfig("dev")

	if cfg.FileDir != "" {
		err := os.Mkdir(cfg.FileDir, 0755)
		if err != nil && !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}
	opts := options.TaskOptions{
		MaxNumResources:       3,
		MaxNumTasks:           3,
		AllowedFileExtensions: cfg.GetAllowedExtensions(),
		FileDir:               cfg.FileDir,
	}

	archiver := &infrastructure.ZipArchiver{}
	httpClient := infrastructure.NewHTTPClient()

	i := interfaces.New(archiver, httpClient)
	s := services.New(i, opts)
	h := handlers.New(s)

	e := handlers.Echo()
	handlers.SetDefault(e)
	handlers.SetAPI(e, h)

	// TODO logging
	// TODO error handling middleware for app errors
	if err := e.Start(":" + strconv.Itoa(cfg.Port)); err != nil {
		panic(err)
	}
}
