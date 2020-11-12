package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/urfave/negroni"
	"net/http"
	"{{ shortName }}/models"
	"path"
	"{{ shortName }}/test_helpers"
	"testing"
)

func insertTestSuperUser(db *sql.DB) {
	db.Exec("INSERT INTO public.members (id, name, email, password, is_superuser) VALUES (1, 'test', 'email', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO', TRUE)")
}

func TestGetMembers(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/api/admin/members", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{
		"members": []interface{}{
			map[string]interface{}{
				"Email":       "email",
				"ID":          "1",
				"IsSuperuser": true,
				"Name":        "test",
			},
		},
		"meta": map[string]interface{}{
			"limit":  float64(20),
			"offset": float64(0),
			"total":  float64(1),
		},
	}, body)
}

func TestGetMembersNonAutorized(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	req, _ := http.NewRequest("GET", "/api/admin/members", nil)
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Expired session or cookie", "Messages": []interface{}{"Session Expired.  Log out and log back in."}}, body)
}

func TestGetMembersForNonSuperuser(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/api/admin/members", nil)
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusForbidden, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{
		"Messages": []interface{}{"Only superuser can get member list"},
		"Status":   "Only superuser can get member list",
	}, body)
}
