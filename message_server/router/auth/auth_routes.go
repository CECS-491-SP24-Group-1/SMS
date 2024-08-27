package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"wraith.me/message_server/config"
	cr "wraith.me/message_server/redis"
	"wraith.me/message_server/schema/user"
)

var (
	// Shared user collection across the entire package.
	uc *user.UserCollection

	// Shared Redis client across the entire package.
	rcl *redis.Client

	// Shared config object across the entire package
	cfg *config.Config

	// Shared env object across the entire package
	env *config.Env
)

// Sets up routes for the `/api/auth` endpoint.
func AuthRoutes(cfgg *config.Config, envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = user.GetCollection() //Init on this is called only once; this is safe to use
	rcl = cr.GetInstance().GetClient()
	cfg = cfgg
	env = envv

	//Add routes (unauthenticated)
	r.Post("/register", RegisterUserRoute)
	r.Post("/login_req", RequestLoginUserRoute)
	r.Post("/login_verify", VerifyLoginUserRoute)
	r.Post("/refresh", RefreshTokenRoute)

	//Add the test route
	authTest := NewAuthTestRouter("", envv)
	r.Group(authTest.Router())

	//Return the router
	return r
}
