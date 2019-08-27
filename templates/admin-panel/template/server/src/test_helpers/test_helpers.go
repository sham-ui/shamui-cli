package test_helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type callsStore map[string]*FunctionCall

type FunctionsCallsStorage struct {
	callsByFunctionName callsStore
}

type fnArgs []interface{}
type fnCallArgs []fnArgs
type FunctionCall struct {
	calls fnCallArgs
}

func (storage *FunctionsCallsStorage) For(name string) *FunctionCall {
	call, ok := storage.callsByFunctionName[name]
	if !ok {
		call = &FunctionCall{
			calls: fnCallArgs{},
		}
		storage.callsByFunctionName[name] = call
	}
	return call
}

func (fc *FunctionCall) Add(args ...interface{}) {
	fc.calls = append(fc.calls, args)
}

func (fc *FunctionCall) Count() int {
	return len(fc.calls)
}

func (fc *FunctionCall) ArgsAt(index int) []interface{} {
	return fc.calls[index]
}

func NewMockFunctionCalls() *FunctionsCallsStorage {
	return &FunctionsCallsStorage{
		callsByFunctionName: make(callsStore),
	}
}

func ClearDB(db *sql.DB) {
	db.Exec("DELETE FROM members")
	db.Exec("ALTER SEQUENCE members_id_seq RESTART WITH 1")
	db.Exec("DELETE FROM http_sessions")
	db.Exec("ALTER SEQUENCE http_sessions_id_seq RESTART WITH 1")
}

func ExecuteRequest(n *negroni.Negroni, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	n.ServeHTTP(rr, req)
	return rr
}

func SetCSRFToken(n *negroni.Negroni, req *http.Request) {
	csrfReq, _ := http.NewRequest("GET", "/api/csrftoken", nil)
	res := ExecuteRequest(n, csrfReq)
	req.Header.Set("Cookie", res.Header().Get("Set-Cookie"))
	req.Header.Set("X-Csrf-Token", res.Header().Get("X-Csrf-Token"))
}

func MergeCookies(req *http.Request, resp *httptest.ResponseRecorder) string {
	const separator = "; "
	merged := map[string]struct{}{}
	cookies := strings.Split(req.Header.Get("Cookie"), separator)
	for _, chunk := range cookies {
		merged[chunk] = struct{}{}
	}
	cookies = strings.Split(resp.Header().Get("Set-Cookie"), separator)
	for _, chunk := range cookies {
		merged[chunk] = struct{}{}
	}
	var chunks []string
	for chunk := range merged {
		chunks = append(chunks, chunk)
	}
	return strings.Join(chunks, separator)
}

func DisableLogger() {
	log.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.FatalLevel)
}

func UnmarshalJSON(data []byte) (map[string]interface{}, error) {
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(data, &jsonData)
	return jsonData, err
}

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
