package users

import (
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/config"
	"wraith.me/message_server/db"
	cr "wraith.me/message_server/redis"
)

var (
	// Shared MongoDB client across the entire package.
	mcl *mongo.Client

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
	mcl = db.GetInstance().GetClient()
	rcl = cr.GetInstance().GetClient()
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
	r.Post("/login_req", RequestLoginUserRoute)
	r.Post("/login_verify", VerifyLoginUserRoute)

	//Return the router
	return r
}
