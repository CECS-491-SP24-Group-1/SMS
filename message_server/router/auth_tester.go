package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/mw"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/util/httpu"
)

//
//-- CLASS: AuthTestRouter
//

// Defines a wrapper around an authentication test router.
type AuthTestRouter struct {
	//The path at which this route is available.
	Path string

	//The token scopes that are valid for the router.
	Scopes []token.TokenScope
}

// Creates a new `AuthTestRouter` object.
func NewAuthTestRouter(path string, scopes []token.TokenScope) AuthTestRouter {
	if path == "" {
		path = "/auth_test"
	}
	return AuthTestRouter{Path: path, Scopes: scopes}
}

// Creates an authentication test route; accessible via a GET request.
func (atr AuthTestRouter) Router() func(r chi.Router) {
	//Set up the response to return if everything goes ok
	successHandler := func(w http.ResponseWriter, r *http.Request) {
		httpu.HttpOkAsJson(w, fmt.Sprintf("authentication succeeded for user %s", r.Header.Get(mw.AuthHttpHeaderSubject)), http.StatusOK)
	}

	//Create the router to respond to the route
	return func(r chi.Router) {
		//Set the auth middleware handler and success responder
		r.Use(mw.NewAuthMiddleware(atr.Scopes))
		r.Get(atr.Path, successHandler)
	}
}
