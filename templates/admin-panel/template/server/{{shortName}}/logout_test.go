package main

import (
	"bytes"
	"encoding/json"
	"github.com/urfave/negroni"
	"net/http"
	"{{ shortName }}/models"
	"path"
	"{{ shortName }}/test_helpers"
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	loginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, loginReq)
	loginResponse := test_helpers.ExecuteRequest(n, loginReq)
	test_helpers.Equals(t, http.StatusOK, loginResponse.Code)

	req, _ := http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)

	req, _ = http.NewRequest("POST", "/api/logout", bytes.NewBuffer([]byte{}))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	test_helpers.Equals(t, "", response.Body.String())

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusUnauthorized, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, "Session is expired.\n", response.Body.String())
}

func TestLogoutFail(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)

	req, _ := http.NewRequest("POST", "/api/logout", bytes.NewBuffer([]byte{}))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	test_helpers.Equals(t, "\"Session Expired.  Log out and log back in.\"\n", response.Body.String())
}
