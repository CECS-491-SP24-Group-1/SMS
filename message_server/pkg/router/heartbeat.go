package router

import (
	"encoding/json"
	"log"
	"net/http"

	"wraith.me/message_server/pkg/db"
	cr "wraith.me/message_server/pkg/redis"
	"wraith.me/message_server/pkg/util"
)

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	//Perform a heartbeat for the Mongo database, Redis, and SMTP server. Then add them together
	dbPing, merr := db.GetInstance().Heartbeat()
	redisPing, rerr := cr.GetInstance().Heartbeat()
	totPing := dbPing + redisPing

	//Create the JSON response
	var payload []byte
	code := http.StatusOK
	if merr == nil && rerr == nil {
		//Construct the response using a map
		resp := map[string]interface{}{
			"status":  "ok",
			"db_ping": totPing,
		}

		//Marshal the map to JSON
		var jerr error
		payload, jerr = json.Marshal(resp)
		if jerr != nil {
			log.Fatalf("couldn't create heartbeat response; %s", jerr)
		}

	} else {
		if merr != nil {
			util.ErrResponse(http.StatusInternalServerError, merr).Respond(w)
		} else if rerr != nil {
			util.ErrResponse(http.StatusInternalServerError, rerr).Respond(w)
		}
		return
	}

	//Respond to the request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
