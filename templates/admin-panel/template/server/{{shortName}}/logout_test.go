package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
	"time"
)

func TestLogoutSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")

	resp = env.API.Request("POST", "/api/logout", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, "", resp.Text(), "text")

	time.Sleep(100 * time.Millisecond) // Session reset in background
	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON(), "body")
}

func TestLogoutFail(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/logout", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON(), "body")
}
