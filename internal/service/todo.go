package service

import (
	"encoding/json"
	"github.com/harsha-aqfer/todo/pkg"
	"net/http"
)

func (s *Service) listTodos(w http.ResponseWriter, r *http.Request) error {
	all := r.URL.Query().Get("all") == "true"

	todos, err := s.store.ListTodos(all)

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return WriteJSON(w, http.StatusOK, todos)
}

func (s *Service) createTodo(w http.ResponseWriter, r *http.Request) error {
	var tr pkg.TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// TODO: add validation

	if err := s.store.CreateTodo(&tr); err != nil {
		return WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return WriteJSON(w, http.StatusOK, nil)
}

func (s *Service) HandleTodos(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.listTodos(w, r)
	case http.MethodPost:
		return s.createTodo(w, r)
	default:
		return WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
