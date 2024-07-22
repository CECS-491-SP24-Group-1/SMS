package challenges

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/config"
	"wraith.me/message_server/schema/user"
)

// Shared user collection across the entire package.
var uc *user.UserCollection

// Shared env object across the entire package
var env *config.Env

// Sets up routes for the `/challenges` endpoint.
func ChallengeRoutes(envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = user.InitCollection()
	env = envv

	//Add routes
	r.Get("/email/{ctext}", SolveEChallengeRoute)

	//Add routes
	//r.Get("/{id}/solve", SolveEChallengeRoute)
	//r.Get("/{id}/status", GetChallengeRoute)

	//Return the router
	return r
}
