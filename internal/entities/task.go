package entities

import (
	"github.com/google/uuid"
	"sync"
)

type Task struct {
	ID          uuid.UUID
	Resources   []Resource
	Status      Status
	ArchivePath string

	Mu *sync.RWMutex
}

func NewTask() *Task {
	return &Task{
		ID:        uuid.New(),
		Status:    StatusPending,
		Resources: make([]Resource, 0),
		Mu:        &sync.RWMutex{},
	}
}
