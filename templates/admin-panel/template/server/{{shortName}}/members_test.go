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

func TestUpdateNameSuccess(t *testing.T) {
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

	payload, _ = json.Marshal(map[string]interface{}{
		"NewName": "edited test name",
	})
	req, _ := http.NewRequest("PUT", "/api/members/name", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "OK", "Messages": []interface{}{"edited test name"}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "edited test name", "Email": "email", "IsSuperuser": false}, body)
}

func TestUpdateNameUnauthtorized(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"NewName": "edited test name",
	})
	req, _ := http.NewRequest("PUT", "/api/members/name", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Expired session or cookie", "Messages": []interface{}{"Session Expired.  Log out and log back in."}}, body)
}

func TestUpdateNameShortName(t *testing.T) {
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

	payload, _ = json.Marshal(map[string]interface{}{
		"NewName": "",
	})
	req, _ := http.NewRequest("PUT", "/api/members/name", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Bad Name", "Messages": []interface{}{"Name must have more than 0 characters."}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)
}

func TestUpdateEmailSuccess(t *testing.T) {
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

	payload, _ = json.Marshal(map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	req, _ := http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "OK", "Messages": []interface{}{"newemail@test.com"}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "newemail@test.com", "IsSuperuser": false}, body)
}

func TestUpdateEmailUnauthtorized(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	payload, _ := json.Marshal(map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "newemail@test.com",
	})
	req, _ := http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, req)
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Expired session or cookie", "Messages": []interface{}{"Session Expired.  Log out and log back in."}}, body)
}

func TestUpdateEmailShort(t *testing.T) {
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

	payload, _ = json.Marshal(map[string]interface{}{
		"NewEmail1": "",
		"NewEmail2": "newemail@test.com",
	})
	req, _ := http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Bad Email", "Messages": []interface{}{"Email must have more than 0 characters."}}, body)

	payload, _ = json.Marshal(map[string]interface{}{
		"NewEmail1": "newemail@test.com",
		"NewEmail2": "",
	})
	req, _ = http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Bad Email", "Messages": []interface{}{"Email must have more than 0 characters."}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)
}

func TestUpdateEmailNotMatch(t *testing.T) {
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

	payload, _ = json.Marshal(map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email2",
	})
	req, _ := http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Bad Email", "Messages": []interface{}{"Emails don't match."}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)
}

func TestUpdateEmailNotUnique(t *testing.T) {
	test_helpers.DisableLogger()
	n := negroni.New()
	startApplication(path.Join("testdata", "config.cfg"), n)
	test_helpers.ClearDB(models.Db)
	insertTestUser(models.Db)
	models.Db.Exec("INSERT INTO public.members (id, name, email, password) VALUES (2, 'test', 'email1', '$2a$14$QMQH3E2UyfIKTFvLfguQPOmai96AncIV.1bLbcd5huTG8gZxNfAyO')")
	payload, _ := json.Marshal(map[string]interface{}{
		"Email":    "email",
		"Password": "password",
	})
	loginReq, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(payload))
	test_helpers.SetCSRFToken(n, loginReq)
	loginResponse := test_helpers.ExecuteRequest(n, loginReq)
	test_helpers.Equals(t, http.StatusOK, loginResponse.Code)

	payload, _ = json.Marshal(map[string]interface{}{
		"NewEmail1": "email1",
		"NewEmail2": "email1",
	})
	req, _ := http.NewRequest("PUT", "/api/members/email", bytes.NewBuffer(payload))
	req.Header.Set("Cookie", test_helpers.MergeCookies(loginReq, loginResponse))
	req.Header.Set("X-CSRF-Token", loginResponse.Header().Get("X-CSRF-Token"))
	response := test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusBadRequest, response.Code)
	body, _ := test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Status": "Fail update email", "Messages": []interface{}{"Fail update email"}}, body)

	req, _ = http.NewRequest("GET", "/api/validsession", nil)
	req.Header.Set("Cookie", loginResponse.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", loginResponse.Header().Get("X-Csrf-Token"))
	response = test_helpers.ExecuteRequest(n, req)
	test_helpers.Equals(t, http.StatusOK, response.Code)
	body, _ = test_helpers.UnmarshalJSON(response.Body.Bytes())
	test_helpers.Equals(t, map[string]interface{}{"Name": "test", "Email": "email", "IsSuperuser": false}, body)
}
