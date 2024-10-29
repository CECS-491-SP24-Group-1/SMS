package wschat

import (
	"github.com/olahol/melody"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/util"
)

// Gets the ID of a room from a Melody session context.
func getRoomID(s *melody.Session) *util.UUID {
	//Get the room's ID
	ctx, ok := s.Request.Context().Value(WSChatCtxObjKey).(Context)
	if !ok {
		return nil
	}
	return &ctx.RoomID
}

// Gets the ID of a connecting user from a Melody session context.
func getUserID(s *melody.Session) *util.UUID {
	//Get the user's ID
	ctx, ok := s.Request.Context().Value(WSChatCtxObjKey).(Context)
	if !ok {
		return nil
	}
	return &ctx.MemberID
}

// Gets the membership info of the room.
func getParticipants(s *melody.Session) chatroom.MembershipList {
	//Get the membership info
	ctx, ok := s.Request.Context().Value(WSChatCtxObjKey).(Context)
	if !ok {
		return nil
	}
	return ctx.Participants
}
