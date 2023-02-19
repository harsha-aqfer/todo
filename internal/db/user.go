package db

import (
	"database/sql"
	"fmt"
	"github.com/harsha-aqfer/todo/pkg"
)

type UserDB interface {
	CreateUser(ui *pkg.UserInfo) (int64, error)
	GetUserId(email string) (int64, error)
}

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserDB {
	return &userStore{db: db}
}

func (us *userStore) CreateUser(ui *pkg.UserInfo) (int64, error) {
	r, err := us.db.Exec("INSERT user SET name = ?, email = ?", ui.Name, ui.Email)
	if err != nil {
		return 0, err
	}
	return r.LastInsertId()
}

func (us *userStore) GetUserId(email string) (int64, error) {
	row := us.db.QueryRow("SELECT id FROM user WHERE email = ?", email)

	var id int64
	err := row.Scan(&id)

	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("no such user: %s", email)
	} else if err != nil {
		return 0, err
	}
	return id, nil
}
