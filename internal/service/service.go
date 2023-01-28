package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/harsha-aqfer/todo/internal/db"
	"log"
	"net/http"
)

type Config struct {
	UserName string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

func NewConfig() *Config {
	return &Config{}
}

type Service struct {
	listenAddr string
	store      db.Store
}

func NewService(listenAddr string, c *Config) (*Service, error) {
	store, err := db.NewDB(c.UserName, c.Password, c.Host, c.Database)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Service{
		listenAddr: listenAddr,
		store:      store,
	}, nil
}

func (s *Service) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/todos", makeHTTPHandleFunc(s.HandleTodos))
	router.HandleFunc("/todos/{id}", makeHTTPHandleFunc(s.HandleTodosById))

	log.Println("JSON API server running on port: ", s.listenAddr)

	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}