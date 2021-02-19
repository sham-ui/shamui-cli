package database

import "database/sql"

type TestDatabase struct {
	DB *sql.DB
}

func (db *TestDatabase) Clear() {
	db.DB.Exec("DELETE FROM members")
	db.DB.Exec("ALTER SEQUENCE members_id_seq RESTART WITH 1")
	db.DB.Exec("DELETE FROM http_sessions")
	db.DB.Exec("ALTER SEQUENCE http_sessions_id_seq RESTART WITH 1")
}

func NewTestDatabase(db *sql.DB) *TestDatabase {
	return &TestDatabase{DB: db}
}
