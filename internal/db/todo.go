package db

import (
	"database/sql"
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
	"strings"
	"time"
)

type TodoDB interface {
	ListTodos(userEmail string, all bool) ([]pkg.TodoResponse, error)
	GetTodo(userEmail string, id int64) (*pkg.TodoResponse, error)
	CreateTodo(userId int64, tr *pkg.TodoRequest) error
	UpdateTodo(userId int64, id int64, tr *pkg.TodoRequest) error
	DeleteTodo(userEmail string, id int64) error
}

type todoStore struct {
	db *sql.DB
}

func NewTodoStore(db *sql.DB) TodoDB {
	return &todoStore{db: db}
}

func (ts *todoStore) ListTodos(userEmail string, all bool) ([]pkg.TodoResponse, error) {
	query := "SELECT T.id, T.task, T.category, T.priority, T.created_at, T.completed_at FROM todo T JOIN user U ON T.user_id = U.id WHERE U.email = ?"

	if !all {
		query += " AND NOT done"
	}
	rows, err := ts.db.Query(query, userEmail)
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

func (ts *todoStore) CreateTodo(userId int64, tr *pkg.TodoRequest) error {
	query := "INSERT todo SET user_id = ?, task = ?"

	params := []interface{}{userId, tr.Task}

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

func (ts *todoStore) GetTodo(userEmail string, id int64) (*pkg.TodoResponse, error) {
	query := "SELECT T.id, T.task, T.category, T.priority, T.created_at, T.completed_at FROM todo T JOIN user U ON T.user_id = U.id WHERE T.id = ? AND U.email = ?"

	rows, err := ts.db.Query(query, id, userEmail)
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

func (ts *todoStore) UpdateTodo(userId int64, id int64, tr *pkg.TodoRequest) error {
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

	params = append(params, id, userId)
	_, err := ts.db.Exec(fmt.Sprintf("UPDATE todo SET %s WHERE id = ? AND user_id = ?", strings.Join(qs, ", ")), params...)
	return err
}

func (ts *todoStore) DeleteTodo(userEmail string, id int64) error {
	_, err := ts.db.Exec("DELETE FROM todo WHERE id = ? AND user_id = (SELECT id FROM user WHERE email = ?)", id, userEmail)
	return err
}
