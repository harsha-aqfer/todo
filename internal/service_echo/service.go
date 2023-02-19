package service_echo

import (
	"fmt"
	"github.com/harsha-aqfer/todo/internal/db"
	"github.com/labstack/echo/v4"
	"log"
)

type Config struct {
	UserName   string `json:"user"`
	Password   string `json:"password"`
	Database   string `json:"database"`
	Host       string `json:"host"`
	ListenAddr string `json:"listen_addr"`
}

func NewConfig() *Config {
	return &Config{}
}

type Service struct {
	listenAddr string
	db         *db.DB
}

func NewService(c *Config) (*Service, error) {
	store, err := db.NewDB(c.UserName, c.Password, c.Host, c.Database)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Service{
		listenAddr: c.ListenAddr,
		db:         store,
	}, nil
}

func (s *Service) Run() {
	e := echo.New()

	// Register app (*App) to be injected into all HTTP handlers.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("service", s)
			return next(c)
		}
	})

	e.POST("/v1/todos", createTodo)
	e.GET("/v1/todos", listTodos)

	e.GET("/v1/todos/:id", getTodo)
	e.PUT("/v1/todos/:id", updateTodo)
	e.DELETE("/v1/todos/:id", deleteTodo)

	e.GET("/v1/version", getVersion)

	log.Println("Server running on port: ", s.listenAddr)
	e.Logger.Fatal(e.Start(s.listenAddr))
}
