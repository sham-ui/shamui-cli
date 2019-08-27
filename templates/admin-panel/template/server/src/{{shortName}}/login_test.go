package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/urfave/negroni"
	"net/http"
	"{{ shortName }}/models"
	"path"
	"test_helpers"
	"testing"
)

func insertTestUser(db *sql.DB) {
	db.Exec("INSERT INTO public.members (id, name, email, password) VALUES (1, 'test', 'email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO')")
}

func TestLoginInvalidCSRF(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte{}))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusForbidden, response.Code)
	test_helpers.Equals(t, "Forbidden - CSRF token invalid\n", response.Body.String())
}

func TestLoginSuccess(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"ID": "1", "Status": "OK", "Name": "test", "Email": "email"}, body)
}

func TestLoginIncorrectPassword(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "incorrectPassword",
	})
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusUnauthorized, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Failed to authenticate", "Messages": []interface{}{"Incorrect username or password"}}, body)
}

func TestLoginIncorrectEmail(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "incorrectemail",
		"Password": "password",
	})
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusUnauthorized, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Failed to authenticate", "Messages": []interface{}{"Incorrect username or password"}}, body)
}
