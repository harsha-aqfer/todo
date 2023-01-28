package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/harsha-aqfer/todo/pkg"
	"net/http"
	"reflect"
	"strconv"
)

func getID(r *http.Request) (int64, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func (s *Service) listTodos(w http.ResponseWriter, r *http.Request) error {
	all := r.URL.Query().Get("all") == "true"

	todos, err := s.store.ListTodos(all)

	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todos)
}

func (s *Service) createTodo(w http.ResponseWriter, r *http.Request) error {
	var tr pkg.TodoRequest

	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		return err
	}

	if err := tr.Validate(); err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := s.store.CreateTodo(&tr); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, nil)
}

func (s *Service) getTodo(w http.ResponseWriter, r *http.Request) error {
	todoId, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	todo, err := s.store.GetTodo(todoId)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todo)
}

func (s *Service) updateTodo(w http.ResponseWriter, r *http.Request) error {
	todoId, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var tr pkg.TodoRequest
	if err = json.NewDecoder(r.Body).Decode(&tr); err != nil {
		return err
	}

	if reflect.ValueOf(tr).IsZero() {
		return WriteJSON(w,
			http.StatusBadRequest,
			map[string]string{"error": "empty body is not supported"},
		)
	}

	if err = s.store.UpdateTodo(todoId, &tr); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, nil)
}

func (s *Service) deleteTodo(w http.ResponseWriter, r *http.Request) error {
	todoId, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err = s.store.DeleteTodo(todoId); err != nil {
		return err
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

func (s *Service) HandleTodosById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.getTodo(w, r)
	case http.MethodPut:
		return s.updateTodo(w, r)
	case http.MethodDelete:
		return s.deleteTodo(w, r)
	default:
		return WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}
