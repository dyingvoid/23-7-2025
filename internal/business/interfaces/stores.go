package interfaces

import (
	"23-7-2025/internal/entities"
	"github.com/google/uuid"
)

type TaskStore interface {
	Set(*entities.Task)
	Get(uuid.UUID) (*entities.Task, bool)
	AddResource(id uuid.UUID, resource entities.Resource)
	Len() int
}
