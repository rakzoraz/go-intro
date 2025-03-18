package handler

import (
	"encoding/json"
	"go-intro/internal/model"
	"net/http"
)

type Handler struct {
	todos map[int]string
}

func (h *Handler) GetTodos(writer http.ResponseWriter, request *http.Request) {
	h.todos = map[int]string{
		1: "neki1",
		2: "neki2",
	}

	response := &model.TodosResponse{}
	if len(h.todos) == 0 {
		response.Message = "no todos in memory"
		writer.WriteHeader(http.StatusNotFound)
	}

	response.Todos = h.todos
	b, err := json.MarshalIndent(response, "", "  ")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(b)

}
