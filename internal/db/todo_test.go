package db

import (
	"database/sql/driver"
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/harsha-aqfer/todo/pkg"
	asserts "github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func Test_Db_ListTodos(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	now := time.Now()
	end := now.Add(20 * time.Minute)

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		rows := sqlMock.NewRows([]string{"id", "task", "category", "priority", "created_at", "completed_at"}).
			AddRow(1, "task-1", "work", "low", now, end).
			AddRow(2, "task-2", "home", "low", now, nil)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, task, category, priority, created_at, completed_at FROM todo WHERE NOT done")).WillReturnRows(rows)

		s := &DB{Sql: db, Todo: NewTodoStore(db)}
		output, err := s.Todo.ListTodos(false)

		if assert.Nil(err) {
			assert.Len(output, 2)
			assert.Equal(pkg.TodoResponse{Id: 1, Task: "task-1", Category: "work", Priority: "low", CreatedAt: &now, CompletedAt: &end}, output[0])
			assert.Equal(pkg.TodoResponse{Id: 2, Task: "task-2", Category: "home", Priority: "low", CreatedAt: &now}, output[1])
		}
	}
}

func Test_Db_CreateTodo(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()

		mock.ExpectExec(regexp.QuoteMeta("INSERT todo SET task = ?, category = ?, priority = ?")).
			WithArgs("task-1", "home", "low").WillReturnResult(driver.ResultNoRows)

		s := &DB{Sql: db, Todo: NewTodoStore(db)}
		err = s.Todo.CreateTodo(&pkg.TodoRequest{
			Task: "task-1", Category: "home", Priority: "low",
		})
		assert.Nil(err)
	}
}

func Test_Db_GetTodo(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	now := time.Now()
	end := now.Add(20 * time.Minute)

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()
		row := sqlMock.NewRows([]string{"id", "task", "category", "priority", "created_at", "completed_at"}).
			AddRow(1, "task-1", "work", "low", now, end)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, task, category, priority, created_at, completed_at FROM todo WHERE id = ?")).WillReturnRows(row)

		s := &DB{Sql: db, Todo: NewTodoStore(db)}
		output, err := s.Todo.GetTodo(1)

		if assert.Nil(err) {
			assert.Equal(&pkg.TodoResponse{Id: 1, Task: "task-1", Category: "work", Priority: "low", CreatedAt: &now, CompletedAt: &end}, output)
		}
	}
}

func Test_Db_DeleteTodo(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM todo WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(driver.ResultNoRows)

		s := &DB{Sql: db, Todo: NewTodoStore(db)}
		err = s.Todo.DeleteTodo(1)
		assert.Nil(err)
	}
}

func Test_Db_UpdateTodo(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE todo SET task = ?, category = ?, priority = ? WHERE id = ?")).
			WithArgs("task-1", "home", "low", 1).WillReturnResult(driver.ResultNoRows)

		s := &DB{Sql: db, Todo: NewTodoStore(db)}
		err = s.Todo.UpdateTodo(1, &pkg.TodoRequest{Task: "task-1", Category: "home", Priority: "low"})
		assert.Nil(err)
	}
}
