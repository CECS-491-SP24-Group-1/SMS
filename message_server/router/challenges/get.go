package challenges

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/mw"
	"wraith.me/message_server/util/httpu"
)

// Handles incoming requests made to `GET /challenges/{id}/get`.
func GetChallengeRoute(w http.ResponseWriter, r *http.Request) {
	//Get the ID of the challenge
	cid := chi.URLParam(r, "id")

	//Return a 400 if the ID is not of the proper format
	if !mongoutil.IsValidUUIDv7(cid) {
		httpu.HttpErrorAsJson(w, fmt.Errorf("incorrect ID format; must be a UUIDv7"), 400)
		return
	}

	//Get the challenge from the database by its ID
	//crud.

	names, _ := mcl.ListDatabaseNames(context.TODO(), bson.M{})

	s := fmt.Sprintf("id: %s; x: %v; cfor: %s; csub: %s; srv_id: %s", cid, names, r.Header.Get(mw.AuthHttpHeaderSubject), r.Header.Get(mw.AuthHttpHeaderScope), env.ID.String())

	w.Write([]byte(s))
}
