package room

import (
	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/globals"
	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/ws/wschat"
)

var (
	// Shared user collection across the entire package.
	uc *user.UserCollection

	rc *chatroom.RoomCollection

	// Shared config object across the entire package.
	cfg *config.Config

	// Shared env object across the entire package.
	env *config.Env

	// Shared Melody WS handler for the entire package.
	mel *wschat.Server
)

// Sets up routes for the `/api/chat/room` endpoint.
func RoomRoutes() chi.Router {
	//Create the router
	r := chi.NewRouter()

	//Set the singletons for the entire package
	uc = globals.UC
	rc = globals.RC
	cfg = globals.Cfg
	env = globals.Env

	//Start up Melody
	mel = wschat.GetInstance()

	//Add routes (unauthenticated)
	// (nada)

	//Add routes (authenticated)
	r.Group(func(r chi.Router) {
		//Apply authentication middleware
		r.Use(mw.NewAuthMiddleware(env))

		//Bind routes
		r.Post("/create", CreateRoomRoute)
		r.Get("/list", GetRoomsRoute)
		r.Get("/{roomID}/members", RoomMembersRoute)
		r.Get("/{roomID}", JoinRoomRoute) //TODO: add `/join`
		r.Post("/{roomID}/leave", LeaveRoomRoute)
		r.Get("/{roomID}/add", AddRoomRoute)
	})

	//Return the router
	return r
}
