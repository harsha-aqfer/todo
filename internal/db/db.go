package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/harsha-aqfer/todo/pkg"
)

type DB struct {
	SQL *sql.DB
}

func NewDB(username, password, host, dbname string) (Store, error) {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbname)
	db, err := sql.Open("mysql", connectString)
	if err == nil {
		return &DB{db}, nil
	}
	return nil, err
}

type Store interface {
	ListTodos(all bool) ([]pkg.TodoResponse, error)
	GetTodo() (*pkg.TodoResponse, error)
	CreateTodo(tr *pkg.TodoRequest) error
	UpdateTodo() error
	DeleteTodo() error
}

func (db *DB) ListTodos(all bool) ([]pkg.TodoResponse, error) {
	query := "SELECT id, task, category, priority, created_at, completed_at FROM todo"

	if !all {
		query += " WHERE NOT done"
	}
	rows, err := db.SQL.Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	todos := make([]pkg.TodoResponse, 0)

	for rows.Next() {
		t := pkg.TodoResponse{}
		var ct sql.NullTime

		err = rows.Scan(&t.Id, &t.Task, &t.Category, &t.Priority, &t.CreatedAt, &ct)

		if err != nil {
			return nil, err
		}
		if ct.Valid {
			t.CompletedAt = &ct.Time
		}
		todos = append(todos, t)
	}
	return todos, nil
}

func (db *DB) CreateTodo(tr *pkg.TodoRequest) error {
	query := "INSERT todo SET task = ?"

	params := []interface{}{tr.Task}

	if tr.Category != "" {
		query += ", category = ?"
		params = append(params, tr.Category)
	}

	if tr.Priority != "" {
		query += ", priority = ?"
		params = append(params, tr.Priority)
	}

	_, err := db.SQL.Exec(query, params...)
	return err
}

func (db *DB) GetTodo() (*pkg.TodoResponse, error) {
	return nil, nil
}

func (db *DB) UpdateTodo() error {
	return nil
}

func (db *DB) DeleteTodo() error {
	return nil
}
