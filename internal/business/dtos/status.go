package dtos

import (
	"23-7-2025/internal/entities"
	"github.com/google/uuid"
)

type State struct {
	Status           string           `json:"status"`
	ArchiveURI       string           `json:"archive_uri,omitempty"`
	ResourceStatuses []ResourceStatus `json:"resource_statuses"`
}

func NewTaskStatus(
	urlBuilder func(taskID uuid.UUID) string,
	t *entities.Task,
) State {
	status := State{
		Status:           t.Status.String(),
		ResourceStatuses: GetTaskResourceStatuses(t),
	}
	if t.Status == entities.StatusArchived {
		status.ArchiveURI = urlBuilder(t.ID)
	}

	return status
}

type ResourceStatus struct {
	URI    string `json:"uri"`
	Status string `json:"status"`
}

func NewResourceStatus(r entities.Resource) ResourceStatus {
	if r.Error != nil {
		return ResourceStatus{
			URI:    r.URI,
			Status: "error",
		}
	}
	if !r.Downloaded {
		return ResourceStatus{
			URI:    r.URI,
			Status: "pending",
		}
	}
	return ResourceStatus{
		URI:    r.URI,
		Status: "success",
	}
}

func GetTaskResourceStatuses(task *entities.Task) []ResourceStatus {
	resourceStatuses := make([]ResourceStatus, 0, len(task.Resources))
	for _, resource := range task.Resources {
		resourceStatuses = append(resourceStatuses, NewResourceStatus(resource))
	}
	return resourceStatuses
}
