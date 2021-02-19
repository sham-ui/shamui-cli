package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestGetServerInfo(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Assert(t, len(resp.JSON()) > 0, "has keys")
}

func TestGetServerInfoNonAutorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON())
}

func TestGetServerInfoForNonSuperuser(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/server-info", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{
		"Messages": []interface{}{"Allowed only for superuser"},
		"Status":   "Forbidden",
	}, resp.JSON())
}
