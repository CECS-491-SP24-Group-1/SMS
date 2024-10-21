package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"wraith.me/message_server/pkg/db/mongoutil"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// CreateChatRoom handles the creation of a new chat room.
func CreateChatRoom(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req struct {
		Participants []string `json:"participants"` // Expecting participant UUIDs
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w) // Using ErrResponse from util
		return
	}

	// Validate the input (at least one participant is required)
	if len(req.Participants) == 0 {
		util.ErrResponse(http.StatusBadRequest, fmt.Errorf("at least one participant is required")).Respond(w)
		return
	}

	// For now, we'll skip fetching the authenticated user from the context

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

	// Create a new chat room
	chatRoom := chatroom.ChatRoom{
		ID:           util.MustNewUUID7(), // Using UUID v7 for unique room IDs
		Participants: participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save the chat room in the database
	err := mongoutil.CreateChatRoom(&chatRoom)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	// Respond with the created chat room using PayloadResponse
	util.PayloadOkResponse("Chat room created successfully", chatRoom).Respond(w)
}
