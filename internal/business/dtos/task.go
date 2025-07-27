package dtos

type Task struct {
	ID     string `json:"id"`
	Status State  `json:"state"`
}
