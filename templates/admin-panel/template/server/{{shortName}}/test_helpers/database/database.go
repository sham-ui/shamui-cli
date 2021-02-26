package database

import "database/sql"

type TestDatabase struct {
	DB *sql.DB
}

func (db *TestDatabase) Clear() {
	db.DB.Exec("DELETE FROM http_sessions")
	db.DB.Exec("DELETE FROM members")
}

func NewTestDatabase(db *sql.DB) *TestDatabase {
	return &TestDatabase{DB: db}
}
