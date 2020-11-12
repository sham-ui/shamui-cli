package sessions

import (
	"encoding/json"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type resSession struct {
	Name        string
	Email       string
	IsSuperuser bool
}

// CsrfToken will generate a CSRF Token
func CsrfToken(w http.ResponseWriter, r *http.Request) {
	log.Info("Generating csrf token")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

// ValidSession checks that the session is valid and can user can make requests
func ValidSession(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	if nil == session {
		log.Info("Session is old, must log out log back in.")
		http.Error(w, "Session is expired.", http.StatusUnauthorized)
	} else {
		msg := &resSession{
			Name:        session.Name,
			Email:       session.Email,
			IsSuperuser: session.IsSuperuser,
		}
		json.NewEncoder(w).Encode(msg)
		log.Info("Session is good.")
		w.WriteHeader(http.StatusOK)
	}
}
