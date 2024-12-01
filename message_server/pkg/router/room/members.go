package room

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/maps"
	"wraith.me/message_server/pkg/db/mongoutil"
	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
	"wraith.me/message_server/pkg/ws/wschat"
)

// Handles incoming requests made to `GET /api/chat/room/{roomID}/members`.
func RoomMembersRoute(w http.ResponseWriter, r *http.Request) {
	//Get the room from the request params
	room := getRoomFromQuery(w, r)
	if room == nil {
		return
	}

	//Get the requestor's info
	requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)

	//Construct an aggregation to get the info for the members of the room
	var userInfo []response.UInfo
	lookups := mongoutil.Slice2BsonA(maps.Keys(room.Participants))
	aggregation := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$in", Value: lookups}}},
		}}},
		user.PublicQuery,
	}

	//Perform the aggregation to get the info for the users in the room
	err := uc.Aggregate(r.Context(), aggregation).All(&userInfo)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Get the room info from the websocket server
	wsRoom := wschat.GetInstance().GetRoom(room.ID)

	//Construct the output room membership array
	membershipInfo := make([]response.RoomMember, len(userInfo))
	for i, uinfo := range userInfo {
		//Get the online status for the current user
		isOnline := false
		if wsRoom != nil {
			isOnline = wsRoom.HasUser(uinfo.ID)
		}

		//Construct the membership object
		membershipInfo[i] = response.RoomMember{
			ID:          uinfo.ID,
			Username:    uinfo.Username,
			DisplayName: uinfo.DisplayName,
			IsMe:        uinfo.ID == requestor.ID,
			Role:        room.Participants[uinfo.ID],
			IsOnline:    isOnline,
		}
	}

	//Return the membership info
	util.PayloadOkResponse("", membershipInfo...).Respond(w)
}
