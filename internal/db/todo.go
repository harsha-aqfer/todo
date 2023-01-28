package db

import (
	"database/sql"
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
	"strings"
	"time"
)

type TodoDB interface {
	ListTodos(all bool) ([]pkg.TodoResponse, error)
	GetTodo(id int64) (*pkg.TodoResponse, error)
	CreateTodo(tr *pkg.TodoRequest) error
	UpdateTodo(id int64, tr *pkg.TodoRequest) error
	DeleteTodo(id int64) error
}

type todoStore struct {
	db *sql.DB
}

func NewTodoStore(db *sql.DB) TodoDB {
	return &todoStore{db: db}
}

func (ts *todoStore) ListTodos(all bool) ([]pkg.TodoResponse, error) {
	query := "SELECT id, task, category, priority, created_at, completed_at FROM todo"

	if !all {
		query += " WHERE NOT done"
	}
	rows, err := ts.db.Query(query)
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

func (ts *todoStore) CreateTodo(tr *pkg.TodoRequest) error {
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

	_, err := ts.db.Exec(query, params...)
	return err
}

func (ts *todoStore) GetTodo(id int64) (*pkg.TodoResponse, error) {
	query := "SELECT id, task, category, priority, created_at, completed_at FROM todo WHERE id = ?"

	rows, err := ts.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

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
		return &t, nil
	}
	return nil, nil
}

func (ts *todoStore) UpdateTodo(id int64, tr *pkg.TodoRequest) error {
	var (
		qs     []string
		params []interface{}
	)

	if tr.Task != "" {
		qs = append(qs, "task = ?")
		params = append(params, tr.Task)
	}

	if tr.Category != "" {
		qs = append(qs, "category = ?")
		params = append(params, tr.Category)
	}

	if tr.Priority != "" {
		qs = append(qs, "priority = ?")
		params = append(params, tr.Priority)
	}

	if tr.Done {
		qs = append(qs, "done = ?")
		params = append(params, int64(1))

		qs = append(qs, "completed_at = ?")
		params = append(params, time.Now().UTC())
	}

	params = append(params, id)
	_, err := ts.db.Exec(fmt.Sprintf("UPDATE todo SET %s WHERE id = ?", strings.Join(qs, ", ")), params...)
	return err
}

func (ts *todoStore) DeleteTodo(id int64) error {
	_, err := ts.db.Exec("DELETE FROM todo WHERE id = ?", id)
	return err
}
