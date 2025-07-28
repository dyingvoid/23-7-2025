package dtos

type Task struct {
	ID    string `json:"id"`
	State State  `json:"state"`
}
