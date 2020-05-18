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

func TestGetSuMemberListBundle(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/dist/su_members_list.bundle.js", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	test_helpers.Equals(t, true, len(response.Body.Bytes()) > 0)
}

func TestGetSuMemberListBundleNonAutorized(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("GET", "/dist/su_members_list.bundle.js", nil)
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusUnauthorized, response.Code)
	test_helpers.Equals(t, "Unauthorized\n", string(response.Body.Bytes()))
}

func TestGetSuMemberListBundleForNonSuperuser(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/dist/su_members_list.bundle.js", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusForbidden, response.Code)
	test_helpers.Equals(t, "Forbidden\n", string(response.Body.Bytes()))
}
