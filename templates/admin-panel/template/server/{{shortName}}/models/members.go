package models

import (
	"database/sql"
	_ "github.com/lib/pq" // github.com/lib/pq
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func CreateMemberStructure() {
	log.Info("Create members table")
	_, err := Db.Exec(`
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
		log.Fatal(err)
	}
}

// NewMember is the struct for the member signup process
type NewMember struct {
	Name, Email, Password, Password2 string
	IsSuperuser                      bool
}

type MemberData struct {
	ID, Name, Email string
	IsSuperuser     bool
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateMember creates the new member record
func CreateMember(m *NewMember) error {
	hashedPw, hashErr := hashPassword(m.Password)
	if hashErr != nil {
		log.Warn("Error hashing password: ", hashErr)
		return hashErr
	}
	_, err := Db.Query("INSERT INTO members(name, email, password, is_superuser) VALUES ($1,$2, $3, $4)", m.Name, m.Email, hashedPw, m.IsSuperuser)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func GetMemberData(email string) *MemberData {
	data := &MemberData{}
	sqlErr := Db.QueryRow("SELECT id, email, name, is_superuser FROM members WHERE email =$1", email).Scan(&data.ID, &data.Email, &data.Name, &data.IsSuperuser)
	if sqlErr == sql.ErrNoRows {
		return data
	}
	if sqlErr != nil {
		log.Println(sqlErr)
	}
	return data
}

// UpdateMemberName uses the member ID to insert a new name
func UpdateMemberName(id string, name string) bool {
	_, sqlErr := Db.Query("UPDATE members SET name = $2 WHERE id = $1", id, name)
	if sqlErr == sql.ErrNoRows {
		name = ""
		return false
	}
	if sqlErr != nil {
		log.Println(sqlErr)
		return false
	}
	return true
}

// UpdateMemberEmail uses the member ID to insert a new email
func UpdateMemberEmail(id string, email string) bool {
	_, sqlErr := Db.Query("UPDATE members SET email = $2 WHERE id = $1", id, email)
	if sqlErr == sql.ErrNoRows {
		return false
	}
	if sqlErr != nil {
		log.Println(sqlErr)
		return false
	}
	return true
}

// UpdateMemberEmail uses the member ID to insert a new password
func UpdateMemberPassword(id string, password string) bool {
	hashedPw, err := hashPassword(password)
	if err != nil {
		log.Warn("Error hashing password: ", err)
		return false
	}
	_, err = Db.Query("UPDATE members SET password = $2 WHERE id = $1", id, hashedPw)
	if err == sql.ErrNoRows {
		return false
	}
	if err != nil {
		log.Warn(err)
		return false
	}
	return true
}

func GetMembers(offset, limit int) []*MemberData {
	var members []*MemberData
	rows, err := Db.Query("SELECT id, email, name, is_superuser FROM members ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	defer rows.Close()
	if nil != err {
		log.Println(err)
		return members
	}
	for rows.Next() {
		data := &MemberData{}
		err := rows.Scan(&data.ID, &data.Email, &data.Name, &data.IsSuperuser)
		if nil != err {
			log.Println(err)
			return members
		}
		members = append(members, data)
	}
	return members
}

func GetMembersCount() int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM members").Scan(&count)
	if nil != err {
		log.Println(err)
	}
	return count
}
