package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestCsrfToken(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	resp := env.API.Request("GET", "/api/csrftoken", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, "", resp.Text())
	asserts.Equals(t, 88, len(resp.Response.Header().Get("X-CSRF-Token")))
}

func TestSessionNotExists(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON())
}

func TestSessionExists(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON())
}

func TestSuperuserSession(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": true}, resp.JSON())
}
