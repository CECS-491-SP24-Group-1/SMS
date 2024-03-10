package mw

import (
	"net/http"
	"strconv"
	"time"
)

const (
	TurnaroundTimeHeader = "X-Processing-Time"
)

// Allows for extra info to be tacked onto a response before its sent.
type delayedWriter struct {
	http.ResponseWriter
	body  string
	begin time.Time
}

// Writes the request and adds extra headers
func (dw *delayedWriter) WriteHeader(statusCode int) {
	w := dw.ResponseWriter
	w.Header().Add(TurnaroundTimeHeader, strconv.FormatInt(time.Since(dw.begin).Nanoseconds(), 10))
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(dw.body))
}

// Writes the total request time of all middlewares to outbound responses.
func ProcessTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Wrap the response writer in a delayed writer
		dw := &delayedWriter{ResponseWriter: w}

		//Record start time
		dw.begin = time.Now()

		//Call the next handler in the chain
		next.ServeHTTP(dw, r)
	})
}
