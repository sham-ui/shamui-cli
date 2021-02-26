package main

import (
	"net/http"
	"{{ shortName }}/test_helpers"
	"{{ shortName }}/test_helpers/asserts"
	"testing"
)

func TestUpdateNameSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "edited test name",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Name updated"}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "edited test name", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdateNameUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "edited test name",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON(), "body")
}

func TestUpdateNameShortName(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/name", map[string]interface{}{
		"NewName": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Name must have more than 0 characters."}}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdateEmailSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Email updated"}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "newemail@test.com", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdateEmailUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON(), "body")
}

func TestUpdateEmailShort(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "",
		"NewEmail2": "newemail@test.com",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Email must have more than 0 characters.", "Emails don't match."}}, resp.JSON(), "body")

	resp = env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Email must have more than 0 characters.", "Emails don't match."}}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdateEmailNotMatch(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email2",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Emails don't match."}}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdateEmailNotUnique(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.DB.DB.Exec("INSERT INTO public.members (id, name, email, password) VALUES (2, 'test', 'email1', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO')")
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/email", map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email1",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Email is already in use."}}, resp.JSON(), "body")

	resp = env.API.Request("GET", "/api/validsession", nil)
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, resp.JSON(), "body")
}

func TestUpdatePasswordSuccess(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Password updated"}, resp.JSON(), "body")

	env.API.ResetCSRF()
	env.API.ResetCookies()
	env.API.GetCSRF()

	resp = env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Incorrect username or password"}}, resp.JSON(), "body")

	resp = env.API.Request("POST", "/api/login", map[string]interface{}{
		"Email":    "email",
		"Password": "newpass",
	})
	asserts.Equals(t, http.StatusOK, resp.Response.Code, "code")
}

func TestUpdatePasswordUnauthtorized(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.API.GetCSRF()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusUnauthorized, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Unauthorized", "Messages": []interface{}{"Session Expired. Log out and log back in."}}, resp.JSON(), "body")
}

func TestUpdatePasswordShort(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "",
		"NewPassword2": "newpass",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Password must have more than 0 characters.", "Passwords don't match."}}, resp.JSON(), "body")

	resp = env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass",
		"NewPassword2": "",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Password must have more than 0 characters.", "Passwords don't match."}}, resp.JSON(), "body")
}

func TestUpdatePasswordNotMatch(t *testing.T) {
	env := test_helpers.NewTestEnv()
	revert := env.Default()
	defer revert()
	env.CreateUser()
	env.API.GetCSRF()
	env.API.Login()

	resp := env.API.Request("PUT", "/api/members/password", map[string]interface{}{
		"NewPassword1": "newpass1",
		"NewPassword2": "newpass2",
	})
	asserts.Equals(t, http.StatusBadRequest, resp.Response.Code, "code")
	asserts.Equals(t, map[string]interface{}{"Status": "Bad Request", "Messages": []interface{}{"Passwords don't match."}}, resp.JSON(), "body")
}
