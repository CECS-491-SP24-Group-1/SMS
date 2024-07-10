package challenges

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/config"
	"wraith.me/message_server/db"
	cr "wraith.me/message_server/redis"
)

// Shared MongoDB client across the entire package.
var mcl *mongo.Client

// Shared Redis client across the entire package.
var rcl *redis.Client

// Shared env object across the entire package
var env *config.Env

// Sets up routes for the `/challenges` endpoint.
func ChallengeRoutes(envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	mcl = db.GetInstance().GetClient()
	rcl = cr.GetInstance().GetClient()
	env = envv

	//Add routes
	r.Get("/{id}/solve", SolveEChallengeRoute)
	r.Get("/{id}/status", GetChallengeRoute)

	//Return the router
	return r
}
