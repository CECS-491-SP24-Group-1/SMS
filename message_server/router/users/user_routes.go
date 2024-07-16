package users

import (
	"text/template"

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

	// Email HTML templates
	emailChallTemplate *template.Template
)

// Sets up routes for the `/users` endpoint.
func UserRoutes(cfgg *config.Config, envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = user.InitCollection() //Init on this is called only once; this is safe to use
	rcl = cr.GetInstance().GetClient()
	cfg = cfgg
	env = envv

	//Read in the email templates
	//TODO: use `go/embed` instead of loading from the FS
	var rerr error
	emailChallTemplate, rerr = template.ParseFiles("./template/registration_email/template.html")

	//Report any errors with the reading process
	if rerr != nil {
		panic(rerr)
	}

	//Add routes
	r.Post("/register", RegisterUserRoute)
	r.Post("/login_req", RequestLoginUserRoute)
	r.Post("/login_verify", VerifyLoginUserRoute)

	//Return the router
	return r
}
