package notification

import (
	"fmt"

	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Constructs a new message notification.
func NewMsgNotif(sender user.User, recipient util.UUID, room chatroom.Room) Notification {
	id := util.MustNewUUID7()
	content := fmt.Sprintf("%s sent you a message in chatroom %s",
		sender.DisplayName, room.ID,
	)
	return newNotifBackend(id, recipient, content, TypeNEWMSG, room.ID.String())
}

// Constructs a new outgoing friend request notification.
func OutgoingFRQNotif(sender user.User, recipient util.UUID) Notification {
	id := util.MustNewUUID7()
	content := fmt.Sprintf("%s <ID: %s> would like to be friends with you",
		sender.Username, sender.ID.String(),
	)
	return newNotifBackend(id, recipient, content, TypeFRQNEW, sender.ID.String())
}

// Constructs a new friend request acceptance notification.
func FRQAcceptNotif(respondent user.User, recipient util.UUID) Notification {
	return frqResponderBackend(respondent, recipient, true)
}

// Constructs a new friend request rejection notification.
func FRQRejectNotif(respondent user.User, recipient util.UUID) Notification {
	return frqResponderBackend(respondent, recipient, false)
}

// Handles creating response friend requests.
func frqResponderBackend(respondent user.User, recipient util.UUID, accepted bool) Notification {
	id := util.MustNewUUID7()
	content := fmt.Sprintf("%s <ID: %s> has %s your friend request",
		respondent.Username, respondent.ID.String(),
		util.If(accepted, "accepted", "rejected"),
	)
	typ := util.If(accepted, TypeFRQACCEPT, TypeFRQREJECT)
	return newNotifBackend(id, recipient, content, typ, respondent.ID.String())
}

// Creates new notification objects.
func newNotifBackend(id, recipient util.UUID, content string, typ Type, context string) Notification {
	return Notification{
		ID:        id,
		Recipient: recipient,
		Content:   content,
		Type:      typ,
		Context:   context,
		Expires:   resolveExpiryTime(id.Time()),
	}
}
