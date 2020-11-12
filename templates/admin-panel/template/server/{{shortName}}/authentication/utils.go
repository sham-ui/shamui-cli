package authentication

import (
	"database/sql"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"{{ shortName }}/config"
)

// Credentials are the user provided email and password
type Credentials struct {
	Email, Password string
}

func getCurPassword(email string) (password string, userPresent bool) {
	db, err := sql.Open("postgres", config.DataBase.GetURL())
	if err != nil {
		log.Error(err)
	}
	sqlErr := db.QueryRow("SELECT password FROM members WHERE email = $1", email).Scan(&password)
	if sqlErr == sql.ErrNoRows {
		userPresent = false
		password = ""
		return
	}
	if sqlErr != nil {
		log.Error(sqlErr)
	}
	userPresent = true
	return
}

func passwordsMatch(c Credentials) bool {
	curPw, userPresent := getCurPassword(c.Email)
	if userPresent != true {
		log.Info("User is not in the database")
		return false
	}
	loginPw := []byte(c.Password)
	hashedPw := []byte(curPw)
	if bcrypt.CompareHashAndPassword(hashedPw, loginPw) != nil {
		log.Info("The passwords do not match")
		return false
	}
	return true
}

// DecodeCredentials decodes the JSON data into a struct containing the email and password.DecodeCredentials
func DecodeCredentials(r *http.Request) (c Credentials) {
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Info("Error decoding credentials >>", err)
	}
	return
}
