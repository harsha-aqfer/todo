package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	Sql     *sql.DB
	Version VersionDB
	Todo    TodoDB
}

func NewDB(username, password, host, dbname string) (*DB, error) {
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbname)
	db, err := sql.Open("mysql", connectString)
	if err == nil {
		return &DB{
			Sql:     db,
			Version: NewVersionStore(db),
			Todo:    NewTodoStore(db),
		}, nil
	}
	return nil, err
}
