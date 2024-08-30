package auth

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"time"

	"wraith.me/message_server/mw"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

// Defines the structure of a session.
type session struct {
	IsCurrent bool      `json:"is_current"` //TODO: add when parent token functionality is added
	Created   time.Time `json:"created"`
	Expires   time.Time `json:"expires"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"string"`
}

// Handles incoming requests made to `GET /api/auth/sessions`.
func SessionsRoute(w http.ResponseWriter, r *http.Request) {
	//Get the user from the auth middleware and the auth token ID
	user := r.Context().Value(mw.AuthCtxUserKey).(user.User)
	ptid := r.Header.Get(mw.AuthAccessParentTokID)

	//Collect the tokens into a map; select attributes are added, but not the whole token
	sessions := make(map[string]session)
	for rtid, tok := range user.Tokens {
		//Decrypt the current refresh token
		dtok, err := token.Decrypt(
			tok.Token, env.SK, env.ID, token.TokenTypeREFRESH,
		)
		if err != nil {
			/*
				util.ErrResponse(http.StatusInternalServerError,
					fmt.Errorf("failed to decrypt refresh tok with id %s: %s", tid, err),
				).Respond(w)
				return
			*/
			//Skip this session since it failed to decrypt
			continue
		}

		//Add the session
		sessions[rtid] = session{
			IsCurrent: subtle.ConstantTimeCompare([]byte(ptid), []byte(rtid)) == 1,
			Created:   dtok.Issued,
			Expires:   dtok.Expiry,
			IP:        dtok.IPAddr.String(),
			UserAgent: dtok.UserAgent,
		}
	}

	//Emit the sessions in a payload response
	s := ""
	if len(sessions) != 1 {
		s = "s"
	}
	util.PayloadOkResponse(
		fmt.Sprintf("found %d session%s", len(sessions), s),
		sessions,
	).Respond(w)
}
