package auth

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/globals"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
)

var (
	// Shared user collection across the entire package.
	uc *user.UserCollection

	// Shared config object across the entire package
	cfg *config.Config

	// Shared env object across the entire package
	env *config.Env
)

// Sets up routes for the `/api/auth` endpoint.
func AuthRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = globals.UC
	cfg = globals.Cfg
	env = globals.Env

	//Add routes (unauthenticated)
	r.Post("/register", RegisterUserRoute)
	r.Post("/login_req", RequestLoginUserRoute)
	r.Post("/login_verify", VerifyLoginUserRoute)
	r.Post("/refresh", RefreshTokenRoute)
	r.Post("/logout", LogoutRoute)

	//Add the test route
	authTest := NewAuthTestRouter("")
	r.Group(authTest.Router())

	//Add routes (authenticated)
	r.Group(func(r chi.Router) {
		r.Use(mw.NewAuthMiddleware(env))
		r.Get("/sessions", SessionsRoute)
	})

	//Return the router
	return r
}
