package notifications

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

	// Shared config object across the entire package.
	cfg *config.Config

	// Shared env object across the entire package.
	env *config.Env
)

// Sets up routes for the `/api/notifications` endpoint.
func NotificationsRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = globals.UC
	cfg = globals.Cfg
	env = globals.Env

	//Add routes (authenticated)
	r.Group(func(r chi.Router) {
		r.Use(mw.NewAuthMiddleware(env))
		//r.Get("/get", e)
		//r.Patch("/read/{nid}", e)
		//r.Patch("/unread/{nid}", e)
		//r.Delete("/remove/{nid}", e)
	})

	//Return the router
	return r
}
