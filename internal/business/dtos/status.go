package dtos

import (
	"23-7-2025/internal/entities"
)

type State struct {
	Status           string           `json:"status"`
	Archive          string           `json:"archive"`
	ResourceStatuses []ResourceStatus `json:"resource_statuses"`
}

func NewTaskStatus(t *entities.Task) State {
	status := State{
		Status:           t.Status.String(),
		ResourceStatuses: GetTaskResourceStatuses(t),
	}
	if t.Status == entities.StatusArchived {
		status.Archive = t.ArchivePath
	}

	return status
}

type ResourceStatus struct {
	URI    string `json:"uri"`
	Status string `json:"status"`
}

func NewResourceStatus(r entities.Resource) ResourceStatus {
	if !r.Downloaded {
		return ResourceStatus{
			URI:    r.URI,
			Status: "pending",
		}
	}
	if r.Error != nil {
		return ResourceStatus{
			URI:    r.URI,
			Status: "error",
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
