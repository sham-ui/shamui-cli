package handlers

import (
	"net/http"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
)

func logoutHandler(ctx *handler.Context, _ interface{}) (interface{}, error) {
	session, _ := ctx.GetSession()
	rawSession := session.GetRawSession()
	// Revoke users authentication
	rawSession.Values["authenticated"] = false
	rawSession.Options.MaxAge = -1
	return nil, rawSession.Save(ctx.Request, ctx.Response)
}

func NewLogoutHandler(sessionsStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(
		logoutHandler,
		handler.WithOnlyForAuthenticated(sessionsStore),
		handler.WithoutSerializeResultToJSON(),
	)
}
