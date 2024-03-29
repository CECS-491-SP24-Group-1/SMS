package users

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/config"
)

// Shared MongoDB client across the entire package.
//var mcl *mongo.Client

// Shared Redis client across the entire package.
//var rcl *redis.Client

// Shared env object across the entire package
var env *config.Env

// Sets up routes for the `/users` endpoint.
func UserRoutes(envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	env = envv

	//Add routes
	r.Post("/register", RegisterUserRoute)

	//Return the router
	return r
}
