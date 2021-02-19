package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestSignupInvalidCSRF(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()

	resp := env.API.Request("POST", "/api/members", nil)
	asserts.Equals(t, http.StatusForbidden, resp.Response.Code)
	asserts.Equals(t, "Forbidden - CSRF token invalid\n", resp.Text())
}

func TestSignupSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Member created"}, resp.JSON())
}

func TestSignupInvalidData(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", []string{})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Messages": nil, "Status": "Bad Request"}, resp.JSON())
}

func TestSignupPasswordMustMatch(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("POST", "/api/members", map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "1password",
		"Password2": "2password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Passwords do not match."}}, resp.JSON())
}

func TestSignupEmailUnique(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	data := map[string]interface{}{
		"Name":      "test",
		"Email":     "email",
		"Password":  "password",
		"Password2": "password",
	}
	resp := env.API.Request("POST", "/api/members", data)
	asserts.Equals(t, http.StatusOK, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Member created"}, resp.JSON())

	resp = env.API.Request("POST", "/api/members", data)
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code)
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Email is already in use."}}, resp.JSON())
}
