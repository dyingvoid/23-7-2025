package services

import (
	"23-7-2025/internal/business/interfaces"
	"23-7-2025/internal/entities"
	"fmt"
)

type ResourceService struct {
	client interfaces.HTTPClienter
}

func NewResourceService(client interfaces.HTTPClienter, filedir string) *ResourceService {
	return &ResourceService{
		client: client,
	}
}

func (rs *ResourceService) DownloadResource(
	resource entities.Resource, filedir string,
) entities.Resource {
	if resource.Downloaded {
		return resource
	}

	filename, err := rs.client.DownloadFile(resource.URI, filedir)
	if err != nil {
		resource.Error = fmt.Errorf("couldn't download resource: %w", err)
		return resource
	}

	resource.Filename = filename
	resource.Downloaded = true
	return resource
}
