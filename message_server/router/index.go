package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"wraith.me/message_server/util"
)

// Responds to requests made at the `/index` route.
func Index(w http.ResponseWriter, r *http.Request) {
	//util.HttpOkSimple(w, fmt.Sprintf("Ok; request id: %s; turnaround: %s\n", w.Header().Get(middleware.RequestIDHeader), w.Header().Get(mw.TurnaroundTimeHeader)))
	//w.Header().Set(middleware.TurnaroundTimeHeader, reqID)
	util.OkResponse(fmt.Sprintf("Ok; request id: %s", w.Header().Get(middleware.RequestIDHeader))).Respond(w)
	//util.HttpOkSimple(w, "Ok")
	//util.HttpOkJson(w, "{\"msg\": \"Ok\"}")
}
