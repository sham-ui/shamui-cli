//Package authentication challenges user credentials and creates or destroys cookie based sessions.
package authentication

import (
	"encoding/json"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"net/http"
	"{{ shortName }}/config"
	"{{ shortName }}/models"
	"{{ shortName }}/sessions"
)

type memberDetails struct {
	Status string
	Name   string
	Email  string
	ID     string
}

type errorMessage struct {
	Status   string
	Messages []string
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Login gets a new session for the user if the credential check passes
func Login(w http.ResponseWriter, r *http.Request) {
	store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
	check(err)
	defer store.Close()
	session, err := store.Get(r, "{{ shortName }}-session")
	check(err)

	// Limit the sessions to 1 24-hour day
	session.Options.MaxAge = 86400 * 1
	session.Options.Domain = config.Server.AllowedDomains[0]
	session.Options.HttpOnly = true

	creds := DecodeCredentials(r)
	// Authenticate based on incoming http request
	if passwordsMatch(creds) != true {
		log.Printf("Bad password for member: %v", creds.Email)
		msg := errorMessage{
			Status:   "Failed to authenticate",
			Messages: []string{"Incorrect username or password"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		//http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		json.NewEncoder(w).Encode(msg)
		return
	}
	// Get the memberID based on the supplied email
	memberID := models.GetMemberID(creds.Email)
	memberName := models.GetMemberName(memberID)
	m := memberDetails{
		Status: "OK",
		ID:     memberID,
		Name:   memberName,
		Email:  creds.Email,
	}

	// Respond with the proper content type and the memberID
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	// Set cookie values and save
	session.Values["authenticated"] = true
	session.Values["name"] = memberName
	session.Values["email"] = creds.Email
	session.Values["id"] = memberID
	if err = session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
	}
	json.NewEncoder(w).Encode(m)
}

// Logout destroys the session
func Logout(w http.ResponseWriter, r *http.Request) {
	if sessions.GetSession(r) == nil {
		json.NewEncoder(w).Encode("Session Expired.  Log out and log back in.")
	}
	store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
	check(err)
	defer store.Close()

	session, err := store.Get(r, "{{ shortName }}-session")
	check(err)
	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)
}
