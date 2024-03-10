package users

import (
	"github.com/go-chi/chi/v5"
)

// Sets up routes for the `/users` endpoint.
func UserRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Add routes
	r.Post("/register", RegisterUserRoute)

	//Return the router
	return r
}
