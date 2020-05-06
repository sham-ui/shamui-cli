// Package sessions resolves session related issues
package sessions

import (
	"github.com/antonlindstrom/pgstore"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"{{ shortName }}/config"
)

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

type Session struct {
	ID          string
	Email       string
	Name        string
	IsSuperuser bool
}

func GetSession(r *http.Request) *Session {
	store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
	check(err)
	defer store.Close()

	session, err := store.Get(r, "{{ shortName }}-session")
	check(err)

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return nil
	}
	id, _ := session.Values["id"].(string)
	name, _ := session.Values["name"].(string)
	email, _ := session.Values["email"].(string)
	isSuperuser, _ := session.Values["is_superuser"].(bool)
	return &Session{
		ID:          id,
		Name:        name,
		Email:       email,
		IsSuperuser: isSuperuser,
	}
}
