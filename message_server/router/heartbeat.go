package router

import (
	"encoding/json"
	"log"
	"net/http"

	"wraith.me/message_server/db"
	"wraith.me/message_server/email"
	credis "wraith.me/message_server/redis"
	"wraith.me/message_server/util/httpu"
)

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	//Perform a heartbeat for the Mongo database, Redis, and SMTP server. Then add them together
	dbPing, merr := db.GetInstance().Heartbeat()
	redisPing, rerr := credis.GetInstance().Heartbeat()
	smtpPing, eerr := email.GetInstance().Heartbeat()
	totPing := dbPing + redisPing + smtpPing

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
			httpu.HttpErrorAsJson(w, merr, http.StatusInternalServerError)
		} else if rerr != nil {
			httpu.HttpErrorAsJson(w, rerr, http.StatusInternalServerError)
		} else if eerr != nil {
			httpu.HttpErrorAsJson(w, eerr, http.StatusInternalServerError)
		}
		return
	}

	//Respond to the request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
