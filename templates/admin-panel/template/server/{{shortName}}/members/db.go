package members

import (
	"database/sql"
	"fmt"
)

func CreateMemberStructure(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS members(
			"id" SERIAL PRIMARY KEY, 
			"name" varchar(100), 
			"email" varchar(100),
			"password" varchar(255),
			"is_superuser" boolean NOT NULL DEFAULT false
		);
		ALTER TABLE ONLY members DROP CONSTRAINT IF EXISTS member_email_unique;
		ALTER TABLE ONLY members ADD CONSTRAINT member_email_unique UNIQUE (email);
	`)
	if err != nil {
		return fmt.Errorf("create table: %s", err)
	}
	return nil
}

func IsUniqueEmail(db *sql.DB, email string) (bool, error) {
	var existingEmail string
	row := db.QueryRow("SELECT email FROM members WHERE email = $1", email)
	err := row.Scan(&existingEmail)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("select email: %s", err)
	}
	return false, nil
}

// CreateMember creates the new member record
func CreateMember(db *sql.DB, m *MemberData) error {
	_, err := db.Query("INSERT INTO members(name, email, password, is_superuser) VALUES ($1,$2, $3, $4)", m.Name, m.Email, m.Password, m.IsSuperuser)
	if nil != err {
		return fmt.Errorf("insert into members: %s", err)
	}
	return nil
}

// UpdateMemberName uses the member ID to update a name
func UpdateMemberName(db *sql.DB, id string, name string) error {
	_, err := db.Query("UPDATE members SET name = $2 WHERE id = $1", id, name)
	return err
}

// UpdateMemberEmail uses the member ID to update a email
func UpdateMemberEmail(db *sql.DB, id string, email string) error {
	_, err := db.Query("UPDATE members SET email = $2 WHERE id = $1", id, email)
	return err
}

// UpdateMemberPassword uses the member ID to update a password
func UpdateMemberPassword(db *sql.DB, id string, password string) error {
	_, err := db.Query("UPDATE members SET password = $2 WHERE id = $1", id, password)
	return err
}
