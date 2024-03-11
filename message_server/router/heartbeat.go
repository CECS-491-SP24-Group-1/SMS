package router

import (
	"encoding/json"
	"log"
	"net/http"

	"wraith.me/message_server/db"
	"wraith.me/message_server/util"
)

func Heartbeat(w http.ResponseWriter, r *http.Request) {
	//Perform a heartbeat for the database
	dbPing, err := db.GetInstance().Heartbeat()

	//Create the JSON response
	var payload []byte
	code := http.StatusOK
	if err == nil {
		//Construct the response using a map
		resp := map[string]interface{}{
			"status":  "ok",
			"db_ping": dbPing,
		}

		//Marshal the map to JSON
		payload, err = json.Marshal(resp)
		if err != nil {
			log.Fatalf("couldn't create heartbeat response; %s", err)
		}

	} else {
		util.HttpErrorAsJson(w, err, http.StatusInternalServerError)
		return
	}

	//Respond to the request
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
