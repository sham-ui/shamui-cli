package main

import (
	"bytes"
	"encoding/json"
	"github.com/urfave/negroni"
	"net/http"
	"{{ shortName }}/models"
	"path"
	"test_helpers"
	"testing"
)

func TestSignupInvalidCSRF(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer([]byte{}))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusForbidden, response.Code)
	test_helpers.Equals(t, "Forbidden - CSRF token invalid\n", response.Body.String())
}

func TestSignupSuccess(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	})
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Member Created", "Errors": nil}, body)
}

func TestSignupInvalidData(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer([]byte{}))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Errors": []interface{}{"Error processing new member data.", "Name must not be empty."}, "Status": "Member Not Created"}, body)
}

func TestSignupPasswordMustMatch(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "1password",
		"Password2": "2password",
	})
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Member Not Created", "Errors": []interface{}{"Passwords do not match."}}, body)
}

func TestSignupEmailUnique(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	StartApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	})
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	test_helpers.ExecuteRequest(n, req)

	req, _ = http.NewRequest("POST", "/api/members", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)

	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Member Not Created", "Errors": []interface{}{"Email is already in use."}}, body)
}
