package db

import (
	"database/sql"
	"fmt"
)

type DB struct {
	SQL *sql.DB
}

func (db *DB) ListTasks() {

}

func (db *DB) GetTask() {

}

func (db *DB) CreateTask() {

}

func (db *DB) UpdateTask() {

}

func (db *DB) DeleteTask() {

}

type Store interface {
	ListTasks() ([]Task, error)
	GetTask() (*Task, error)
	DeleteTask() error
	CreateTask() error
	UpdateTask() error
}

func NewDB(username, password, host, dbname string) (Store, err error) {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbname)
	db, err := sql.Open("mysql", connectString)
	if err == nil {
		return &DB{db}, nil
	}
	return nil, err
}
