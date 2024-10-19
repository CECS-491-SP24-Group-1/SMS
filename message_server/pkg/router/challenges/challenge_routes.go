package challenges

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/globals"
	"wraith.me/message_server/pkg/schema/user"
)

var (
	// Shared user collection across the entire package.
	uc *user.UserCollection

	// Shared env object across the entire package
	env *config.Env
)

// Sets up routes for the `/api/challenges` endpoint.
func ChallengeRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = globals.UC
	env = globals.Env

	//Add routes
	r.Get("/email/{ctext}", SolveEChallengeRoute)

	//Add routes
	//r.Get("/{id}/solve", SolveEChallengeRoute)
	//r.Get("/{id}/status", GetChallengeRoute)

	//Return the router
	return r
}
