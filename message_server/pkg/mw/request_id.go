package mw

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func SendRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			w.Header().Set(middleware.RequestIDHeader, reqID)
		}
		next.ServeHTTP(w, r)
	})
}
