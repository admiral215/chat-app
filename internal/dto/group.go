package dto

type GroupCreate struct {
	Name    string   `json:"name" validate:"required,max=100"`
	Members []string `json:"members"`
}
