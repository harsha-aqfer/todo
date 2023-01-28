package db

import "database/sql"

type VersionDB interface {
	GetVersion() (string, error)
}

type versionStore struct {
	db *sql.DB
}

func NewVersionStore(db *sql.DB) VersionDB {
	return &versionStore{db: db}
}

func (v *versionStore) GetVersion() (version string, err error) {
	err = v.db.QueryRow("SELECT VERSION()").Scan(&version)
	return
}
