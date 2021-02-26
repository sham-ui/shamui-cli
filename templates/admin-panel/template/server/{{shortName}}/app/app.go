package app

import (
	"database/sql"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"{{shortName}}/assets"
	authenticationHandlers "{{shortName}}/authentication/handlers"
	"{{shortName}}/config"
	"{{shortName}}/core/database"
	"{{shortName}}/core/sessions"
	"{{shortName}}/members"
	membersHandlers "{{shortName}}/members/handlers"
	serverHandlers "{{shortName}}/server/handlers"
	sessionHandlers "{{shortName}}/session/handlers"
	"strings"
	"time"
)

func StartApplication(configPath string, n *negroni.Negroni) *sql.DB {
	config.LoadConfiguration(configPath)

	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		log.Fatalf("Fail connect to db: %s", err)
	}

	sessionsStore, err := sessions.NewStore(db, config.Session.Secret)
	if nil == err {
		log.Info("Create pg session store")
	} else {
		log.WithError(err).Fatal("Fail create pg session store")
	}
	// Run a background goroutine to clean up expired sessions from the database.
	defer sessionsStore.StopCleanup(sessionsStore.Cleanup(time.Minute * 5))

	err = members.CreateMemberStructure(db)
	if nil != err {
		log.Fatalf("Fail create members table: %s", err)
	} else {
		log.Info("Create members table")
	}

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
	r.HandleFunc("/api/csrftoken", sessionHandlers.NewCsrfTokenHandler()).Methods("GET")
	r.HandleFunc("/api/login", authenticationHandlers.NewLoginHandler(sessionsStore, db)).Methods("POST")
	r.HandleFunc("/api/logout", authenticationHandlers.NewLogoutHandler(sessionsStore)).Methods("POST")

	// Session Routes
	r.HandleFunc("/api/validsession", sessionHandlers.NewValidSessionHandler(sessionsStore)).Methods("GET")

	// Member CRUD routes
	r.HandleFunc("/api/members", membersHandlers.NewSignupHandler(db, sessionsStore)).Methods("POST")
	r.HandleFunc("/api/members/email", membersHandlers.NewUpdateEmailHandler(db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/members/name", membersHandlers.NewUpdateNameHandler(db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/members/password", membersHandlers.NewUpdatePasswordHandler(db, sessionsStore)).Methods("PUT")

	// Superuser sections
	r.HandleFunc("/api/admin/members", membersHandlers.NewListHandler(db, sessionsStore)).Methods("GET")
	r.HandleFunc("/api/admin/members", membersHandlers.NewCreateHandler(db, sessionsStore)).Methods("POST")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewUpdateHandler(db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}/password", membersHandlers.NewResetPasswordHandler(db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewDeleteHandler(db, sessionsStore)).Methods("DELETE")
	r.HandleFunc("/api/admin/server-info", serverHandlers.NewInfoHandler(sessionsStore)).Methods("GET")

	// Resources
	spaHandler := assets.NewHandler(sessionsStore)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(CSRF(r))

	return db
}
