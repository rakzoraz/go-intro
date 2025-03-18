package model

type TodosResponse struct {
	Message string         `json:"message,omitempty"`
	Todos   map[int]string `json:"todos,omitempty"`
	Error   string         `json:"error,omitempty"`
}
