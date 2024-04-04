package users

import (
	"text/template"

	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/config"
)

// Shared MongoDB client across the entire package.
//var mcl *mongo.Client

// Shared Redis client across the entire package.
//var rcl *redis.Client

// Shared config object across the entire package
var cfg *config.Config

// Shared env object across the entire package
var env *config.Env

// Email HTML templates
var emailChallTemplate *template.Template

// Sets up routes for the `/users` endpoint.
func UserRoutes(cfgg *config.Config, envv *config.Env) chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	cfg = cfgg
	env = envv

	//Read in the email templates
	var rerr error
	emailChallTemplate, rerr = template.ParseFiles("./template/registration_email/template.html")

	//Report any errors with the reading process
	if rerr != nil {
		panic(rerr)
	}

	//Add routes
	r.Post("/register", RegisterUserRoute)

	//Return the router
	return r
}
