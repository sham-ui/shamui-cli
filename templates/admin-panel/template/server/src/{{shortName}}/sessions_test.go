package main

import (
	"net/http"
	"{{ shortName }}/models"
	"path"
	"test_helpers"
	"testing"
	"github.com/urfave/negroni"
	"encoding/json"
	"bytes"
)

func TestCsrfToken(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("GET", "/api/csrftoken", nil)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	test_helpers.Equals(t, "", response.Body.String())
	test_helpers.Equals(t, 88, len(response.Header().Get("X-CSRF-Token")))
}

func TestSessionNotExists(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("GET", "/api/validsession", nil)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusUnauthorized, response.Code)
	test_helpers.Equals(t, "Session is expired.\n", response.Body.String())
}

func TestSessionExists(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", response.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", response.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)
}

func TestSuperuserSession(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	models.Db.Exec("INSERT INTO public.members (id, name, email, password, is_superuser) VALUES (1, 'super_test', 'super_email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO', true)")
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "super_email",
		"Password": "password",
	})
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", response.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", response.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "super_test", "Email": "super_email", "IsSuperuser": true}, body)
}
