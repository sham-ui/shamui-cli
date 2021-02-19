package main

import (
	"net/http"
	"{{shortName}}/test_helpers"
	"{{shortName}}/test_helpers/asserts"
	"testing"
)

func TestGetMembers(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateSuperUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/members?limit=20", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{
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
	}, resp.JSON())
}

func TestGetMembersNonAutorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("GET", "/api/admin/members", nil)
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON())
}

func TestGetMembersForNonSuperuser(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("GET", "/api/admin/members", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{
		"Messages": []interface{}{"Allowed only for superuser"},
		"Status":   "Forbidden",
	}, resp.JSON())
}
