package wschat

import "wraith.me/message_server/pkg/obj"

var (
	// The context key name for a ws chat room ID.
	WSChatCtxObjKey = obj.CtxKey{S: "roomID"}
)
