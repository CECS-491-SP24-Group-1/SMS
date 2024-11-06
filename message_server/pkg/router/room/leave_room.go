package room

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/mw"
	chatroom "wraith.me/message_server/pkg/schema/chat_room"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// LeaveRoomRoute handles requests to `POST /api/chat/room/{roomID}/leave`.
func LeaveRoomRoute(w http.ResponseWriter, r *http.Request) {
    // Extract room ID from URL
    roomID := chi.URLParam(r, "roomID")
    rid, err := util.ParseUUIDv7(roomID)
    if err != nil {
        util.ErrResponse(http.StatusBadRequest, fmt.Errorf("bad room ID format; it must be a UUIDv7")).Respond(w)
        return
    }

    // Get the current user from the context
    requestor := r.Context().Value(mw.AuthCtxUserKey).(user.User)

    // Retrieve the room from the database
    var room chatroom.Room
    err = rc.FindID(r.Context(), rid).One(&room)
    if err != nil {
        code := http.StatusInternalServerError
        if qmgo.IsErrNoDocuments(err) {
            code = http.StatusNotFound
            err = fmt.Errorf("chat room with ID %s not found", rid)
        }
        util.ErrResponse(code, err).Respond(w)
        return
    }

    // Check if the user is a participant in the room
    if _, exists := room.Participants[requestor.ID]; !exists {
        util.ErrResponse(http.StatusForbidden, fmt.Errorf("you are not a member of this room")).Respond(w)
        return
    }

    // Remove the user from the participants list
    delete(room.Participants, requestor.ID)

    // Update the room in the database
    filter := bson.M{"_id": rid}
    update := bson.M{"$set": bson.M{"participants": room.Participants}}
    err = rc.UpdateOne(context.Background(), filter, update)
    if err != nil {
        util.ErrResponse(http.StatusInternalServerError, fmt.Errorf("failed to update room")).Respond(w)
        return
    }

    // test message
    util.PayloadOkResponse("Successfully left the chat room", struct{}{}).Respond(w)


}
