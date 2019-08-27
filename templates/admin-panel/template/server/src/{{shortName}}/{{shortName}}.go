package main

import (
	"assets"
	"github.com/antonlindstrom/pgstore"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"{{ shortName }}/authentication"
	"{{ shortName }}/config"
	"{{ shortName }}/members"
	"{{ shortName }}/models"
	"{{ shortName }}/sessions"
	"strconv"
	"strings"
	"time"
)

func StartApplication(configPath string, n *negroni.Negroni) {
	config.LoadConfiguration(configPath)

	store, err := pgstore.NewPGStore(config.DataBase.GetURL(), []byte(config.Session.Secret))
	if err == nil {
		log.Info("Create pg session store")
	} else {
		log.WithError(err).Fatal("Fail create pg session store")
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	models.ConnToDB(config.DataBase.GetURL())
	models.CreateMemberStructure()

	log.Infof("Allowed domains: %s", strings.Join(config.Server.AllowedDomains, ", "))

	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("{{ shortName }}_csrf"),
		csrf.Secure(false), // Disabled for localhost non-https debugging
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Server.AllowedDomains,
		AllowedMethods:   []string{"PUT", "POST", "GET"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	r := mux.NewRouter()

	// Authentication Routes
	r.HandleFunc("/api/csrftoken", sessions.CsrfToken).Methods("GET")
	r.HandleFunc("/api/login", authentication.Login).Methods("POST")
	r.HandleFunc("/api/logout", authentication.Logout).Methods("POST")
	// Session Routes
	r.HandleFunc("/api/validsession", sessions.ValidSession).Methods("GET")

	// Member CRUD routes
	r.HandleFunc("/api/members", members.SignupMember).Methods("POST")
	r.HandleFunc("/api/members/email", members.UpdateMemberEmail).Methods("PUT")
	r.HandleFunc("/api/members/name", members.UpdateMemberName).Methods("PUT")

	// Resources
	spaHandler := assets.NewHandler()
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(CSRF(r))
}

func main() {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	StartApplication("config.cfg", n)
	port := strconv.Itoa(config.Server.Port)
	log.Infof("Server start on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, context.ClearHandler(n)))
}
