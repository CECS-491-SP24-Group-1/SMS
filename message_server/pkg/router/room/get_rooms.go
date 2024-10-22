package room

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `GET /api/chat/room/list`.
func GetRoomsRoute(w http.ResponseWriter, r *http.Request) {
	//Get the requestor's info
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User) //This assert is safe

	//Construct the search query
	innerKey := "participants." + requestor.ID.String()
	query := bson.D{{Key: innerKey, Value: bson.D{{Key: "$exists", Value: true}}}}

	//Search for the requestor in the rooms collection
	rooms := make([]chatroom.Room, 0)
	err := rc.Find(r.Context(), query).All(&rooms)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//fmt.Printf("rooms: %+v\n", rooms)
	//Return the room list to the user
	util.PayloadOkResponse("", rooms...).Respond(w)
}
