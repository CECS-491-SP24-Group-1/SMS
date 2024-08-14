package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/config"
	"wraith.me/message_server/mw"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

//
//-- CLASS: AuthTestRouter
//

// Defines a wrapper around an authentication test router.
type AuthTestRouter struct {
	//The path at which this route is available.
	Path string

	//The token scopes that are valid for the router.
	//Scopes []token.TokenScope

	//The secrets of the server including ID and encryption key.
	secrets *config.Env
}

// Creates a new `AuthTestRouter` object.
func NewAuthTestRouter(path string, secrets *config.Env) AuthTestRouter {
	if path == "" {
		path = "/test"
	}
	return AuthTestRouter{Path: path, secrets: secrets}
}

// Creates an authentication test route; accessible via a GET request.
func (atr AuthTestRouter) Router() func(r chi.Router) {
	//Set up the response to return if everything goes ok
	successHandler := func(w http.ResponseWriter, r *http.Request) {
		//Get the user object from the auth middleware
		user := r.Context().Value(mw.AuthCtxUserKey).(user.User)

		//Marshal to a map using mapstructure
		ms := make(map[string]interface{})
		if err := util.MSTextMarshal(user, &ms, "bson"); err != nil {
			util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
			return
		}

		//Redact some fields as a test
		ms["id"] = ms["UUID"]
		delete(ms, "UUID")
		delete(ms, "flags")
		delete(ms, "last_ip")
		delete(ms, "options")
		delete(ms, "tokens")

		//Marshal the map to JSON
		jsons, err := json.Marshal(&ms)
		if err != nil {
			util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		}

		fmt.Printf("jsons: `%s\n", jsons)

		//Respond to the user
		util.PayloadOkResponse(
			fmt.Sprintf("authentication succeeded for user with ID %s", user.ID),
			string(jsons),
		).Respond(w)
	}

	//Create the router to respond to the route
	return func(r chi.Router) {
		//Set the auth middleware handler and success responder
		r.Use(mw.NewAuthMiddleware(atr.secrets))
		r.Get(atr.Path, successHandler)
	}
}
