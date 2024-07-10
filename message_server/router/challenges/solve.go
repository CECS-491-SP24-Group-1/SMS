package challenges

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"wraith.me/message_server/crud"
	ch "wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/util"
	"wraith.me/message_server/util/httpu"
)

/*
Handles incoming requests made to `GET /challenges/{id}/solve`. This route
is only to be used for solving email-based challenges.
*/
func SolveEChallengeRoute(w http.ResponseWriter, r *http.Request) {
	//Get the ID of the challenge
	cid := chi.URLParam(r, "id")

	//Return a 400 if the ID is not of the proper format
	if !util.IsValidUUIDv7(cid) {
		httpu.HttpErrorAsJson(w, fmt.Errorf("incorrect ID format; must be a UUIDv7"), http.StatusBadRequest)
		return
	}

	//Attempt to get the challenge from the database
	challs, err := crud.GetChallengesById(mcl, rcl, r.Context(), util.UUIDFromString(cid))
	if err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusInternalServerError)
		return
	}
	chall := challs[0]

	//If the scope is not email, then reject the solve; this route is for solving email challenges only
	if chall.Scope != ch.ChallengeScopeEMAIL {
		httpu.HttpErrorAsJson(w, fmt.Errorf("GET requests to the challenge solve route are for email-based challenges only"), http.StatusBadRequest)
		return
	}

	solve := fmt.Sprintf("solve challenge with id %s, scope %s", cid, chall.Scope)

	//names, _ := mcl.ListDatabaseNames(context.TODO(), bson.M{})

	//s := fmt.Sprintf("id: %s; x: %v; cfor: %s; csub: %s; srv_id: %s", cid, names, r.Header.Get(mw.AuthHttpHeaderSubject), r.Header.Get(mw.AuthHttpHeaderScope), env.ID.String())

	w.Write([]byte(solve))
}
