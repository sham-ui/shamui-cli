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

func TestGetServerInfo(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestSuperUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	loginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, loginReq)
	loginResponse := test_helpers.ExecuteRequest(n, loginReq)
	test_helpers.Equals(t, http.StatusOK, loginResponse.Code)

	req, _ := http.NewRequest("GET", "/api/admin/server-info", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	_, err := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Ok(t, err)
}

func TestGetServerInfoNonAutorized(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("GET", "/api/admin/server-info", nil)
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Expired session or cookie", "Messages": []interface{}{"Session Expired.  Log out and log back in."}}, body)
}

func TestGetServerInfoForNonSuperuser(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/api/admin/server-info", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusForbidden, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{
		"Messages": []interface{}{"Only superuser can get server info"},
		"Status":   "Only superuser can get server info",
	}, body)
}
