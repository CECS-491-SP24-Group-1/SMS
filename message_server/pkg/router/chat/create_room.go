package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wraith.me/message_server/pkg/schema/user"

	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `POST /api/chat/room/create`.
func CreateRoomRoute(w http.ResponseWriter, r *http.Request) {
	//Parse the request body
	var req struct {
		Participants []util.UUID `json:"participants"` //TODO: eventually make a concrete object
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w) // Using ErrResponse from util
		return
	}

	//Validate the input (at least one participant is required)
	if len(req.Participants) == 0 {
		util.ErrResponse(http.StatusBadRequest, fmt.Errorf("at least one participant is required")).Respond(w)
		return
	}

	// For now, we'll skip fetching the authenticated user from the context

	/*
		// Convert participant strings to util.UUID
		var participants []util.UUID
		for _, p := range req.Participants {
			participantUUID := util.UUIDFromString(p)
			if participantUUID.IsNil() {
				util.ErrResponse(http.StatusBadRequest, fmt.Errorf("invalid participant ID format")).Respond(w)
				return
			}
			participants = append(participants, participantUUID)
		}
	*/

	//Get the requestor's info
	owner := r.Context().Value(mw.AuthCtxUserKey).(user.User) //This assert is safe

	//Create a new chat room
	room := chatroom.NewRoom(owner.ID, req.Participants...)

	//Save the chat room in the database
	_, err := chatroom.GetCollection().InsertOne(r.Context(), room)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Respond with the created chat room using PayloadResponse
	util.PayloadOkResponse("Chat room created successfully", room).Respond(w)
}
