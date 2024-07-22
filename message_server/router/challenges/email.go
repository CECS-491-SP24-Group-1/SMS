package challenges

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

/*
Handles incoming requests made to `GET /challenges/email/{ctext}`. This route
is only to be used for solving email-based challenges.
*/
func SolveEChallengeRoute(w http.ResponseWriter, r *http.Request) {
	//Get the challenge text
	ctext := chi.URLParam(r, "ctext")
	if strings.TrimSpace(ctext) == "" {
		//Bail out if nothing was supplied
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("received empty challenge response"),
		).Respond(w)
		return
	}

	//Attempt to decrypt the challenge
	//From this point on, it's safe to assume the user successfully passed the challenge
	ctoken, err := challenge.Decrypt(ctext, env.SK, env.ID, challenge.CPurposeCONFIRM)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return
	}

	//Get the user mentioned in the challenge from the database
	var user user.User
	err = uc.Find(r.Context(), bson.M{"_id": ctoken.SubjectID}).One(&user)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return
	}

	//Mark the user's email as verified and upsert the user into the collection
	user.MarkEmailVerified()
	_, err = uc.Upsert(r.Context(), bson.M{"_id": user.ID}, &user)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return
	}

	//Return the status of the verification
	msg := fmt.Sprintf("email %s successfully verified for user with ID %s", ctoken.Claim, user.ID)
	util.OkResponse(msg).Respond(w)
}
