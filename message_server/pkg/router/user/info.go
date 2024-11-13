package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `GET /api/user/{uid}`.
func HandleInfoRoute(w http.ResponseWriter, r *http.Request) {
	//Extract user ID or username from the URL
	userId := chi.URLParam(r, "uid")

	//Initialize an inner find query for MongoDB; start out with an ID
	var query bson.E

	//If the info string is not a UUID, query by username instead
	uid, err := util.ParseUUIDv7(userId)
	validUUID := err == nil
	if !validUUID {
		query = bson.E{Key: "username", Value: userId}
	} else {
		query = bson.E{Key: "_id", Value: uid}
	}

	//Run the query
	var user user.User
	err = uc.Find(r.Context(), bson.D{query}).One(&user)

	//Check if something went wrong during the query
	if err != nil {
		//Check if the error has to do with a lack of documents
		code := http.StatusInternalServerError
		desc := err.Error()
		if errors.Is(err, mongo.ErrNoDocuments) {
			//Change the error to be a 404
			code = http.StatusNotFound
			desc = fmt.Sprintf(
				"No such user exists by %s %s",
				util.If(validUUID, "UUID", "username"),
				userId,
			)
		}

		//Respond back with an error
		util.ErrResponse(code, errors.New(desc)).Respond(w)
		return
	}

	//Respond back with the user's info
	sendUserInfo(w, user)
}

// Handles incoming requests made to `GET /api/user/me` and `GET /api/user/`.
func HandleMyInfoRoute(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)
	sendUserInfo(w, user)
}

// Sends user info to an HTTP response; truncated.
func sendUserInfo(w http.ResponseWriter, user user.User) {
	info := response.FromUser(user)
	util.PayloadOkResponse("", info).Respond(w)
}
