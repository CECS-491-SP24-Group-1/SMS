package chat

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/globals"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/ws/wschat"
)

var (
	// Shared user collection across the entire package.
	uc *user.UserCollection

	// Shared config object across the entire package.
	cfg *config.Config

	// Shared env object across the entire package.
	env *config.Env

	// Shared Melody WS handler for the entire package.
	mel *wschat.Server
)

// Sets up routes for the `/api/chat` endpoint.
func ChatRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = globals.UC
	cfg = globals.Cfg
	env = globals.Env

	//Start up Melody
	mel = wschat.GetInstance()

	//Add routes (unauthenticated)
	r.Get("/room/{roomID}", JoinRoomRoute)

	//Add routes (authenticated)
	/*
		r.Group(func(r chi.Router) {
			//r.Use(mw.NewAuthMiddleware(env))
			//r.Get("/sessions", SessionsRoute)
		})
	*/

	//Add routes (authenticated)
	r.Group(func(r chi.Router) {
		//Apply authentication middleware
		r.Use(mw.NewAuthMiddleware(env))

		//Route to create a new chat room
		r.Post("/room/create", CreateRoomRoute)
	})

	//Return the router
	return r
}
