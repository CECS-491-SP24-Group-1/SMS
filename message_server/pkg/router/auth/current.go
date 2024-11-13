package auth

import (
	"fmt"
	"net/http"

	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/mw"
	"wraith.me/message_server/pkg/obj/token"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `GET /api/auth/current`.
func CurrentSeshRoute(w http.ResponseWriter, r *http.Request) {
	//Get the current user's info and the access token used
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)
	tok := r.Context().Value(mw.AuthCtxAccessTokKey).(token.Token)

	//Get the parent refresh token from the user
	//It is assumed that the token does exist since auth doesn't pass orphaned access tokens
	ertok := user.Tokens[tok.Parent.String()]

	//Decrypt the refresh token
	rtok, err := token.Decrypt(
		ertok.Token, env.SK, env.ID, token.TokenTypeREFRESH,
	)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError,
			fmt.Errorf(
				"failed to decrypt refresh tok with id %s: %s",
				tok.Parent.String(), err,
			),
		).Respond(w)
		return
	}

	//Construct the output objects
	parTok := response.Session{
		ID:        rtok.ID.String(),
		IsCurrent: true,
		Created:   rtok.Issued,
		Expires:   rtok.Expiry,
		IP:        rtok.IPAddr.String(),
		UserAgent: rtok.UserAgent,
	}
	out := response.AccessSession{
		ID:      tok.ID.String(),
		Created: tok.Issued,
		Expires: tok.Expiry,
		Parent:  parTok,
	}

	//Respond back with the session info
	util.PayloadOkResponse("", out).Respond(w)
}
