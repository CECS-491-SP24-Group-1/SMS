package httpu

import (
	"encoding/json"
	"fmt"
	"net/http"

	"wraith.me/message_server/obj"
)

// Deprecated: Use `util.HttpResponse` instead.
// Writes plaintext to an HTTP response and issues a `500` status code.
func HttpErrorSimple(w http.ResponseWriter, msg string) {
	TextualResponse(w, []byte(msg), http.StatusInternalServerError, "", "")
}

// Deprecated: Use `util.HttpResponse` instead.
// Writes JSON to a request and issues a `500` status code.
func HttpErrorJson(w http.ResponseWriter, json string) {
	TextualResponse(w, []byte(json), http.StatusInternalServerError, "application/json", "")
}

// Deprecated: Use `util.HttpResponse` instead.
// Writes plaintext to an HTTP response and issues a `200` status code.
func HttpOkSimple(w http.ResponseWriter, msg string) {
	TextualResponse(w, []byte(msg), http.StatusOK, "", "")
}

// Deprecated: Use `util.HttpResponse` instead.
// Writes JSON to a request and issues a `200` status code.
func HttpOkJson(w http.ResponseWriter, json string) {
	TextualResponse(w, []byte(json), http.StatusOK, "application/json", "")
}

// Deprecated: Use `util.HttpResponse` instead.
// Issues an HTTP error response as JSON.
func HttpErrorAsJson(w http.ResponseWriter, err error, code int) {
	//Ensure the code is valid
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	//Construct an error response object
	resp := obj.Response{
		Code:   code,
		Status: "error",
		Desc:   err.Error(),
	}

	//Set the mimetype of the output to be JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Write the object to the response writer
	w.WriteHeader(code)
	jerr := json.NewEncoder(w).Encode(&resp)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
	}
}

// Deprecated: Use `util.HttpResponse` instead.
// Issues an HTTP error response as JSON.
func HttpMultipleErrorsAsJson(w http.ResponseWriter, errs []error, code int) {
	//Ensure the code is valid
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	//Populate a list of stringified errors
	stringErrs := make([]string, len(errs))
	for i, err := range errs {
		stringErrs[i] = err.Error()
	}

	//Construct an error response object
	resp := obj.MultiResponse{
		Code:   code,
		Status: "multiple_errors",
		Desc:   stringErrs,
	}

	//Set the mimetype of the output to be JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Write the object to the response writer
	w.WriteHeader(code)
	jerr := json.NewEncoder(w).Encode(&resp)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
	}
}

// Deprecated: Use `util.HttpResponse` instead.
// Issues an HTTP info response as JSON.
func HttpOkAsJson(w http.ResponseWriter, msg string, code int) {
	//Ensure the code is valid
	if code <= 0 {
		code = http.StatusOK
	}

	//Construct an error response object
	resp := obj.Response{
		Code:   code,
		Status: "ok",
		Desc:   msg,
	}

	//Set the mimetype of the output to be JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Write the object to the response writer
	w.WriteHeader(code)
	jerr := json.NewEncoder(w).Encode(&resp)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
	}
}

// Deprecated: Use `util.HttpResponse` instead.
// Utility to write response text to an HTTP response object.
func TextualResponse(w http.ResponseWriter, payload []byte, code int, mime string, encoding string) {
	if mime == "" {
		mime = "text/plain"
	}
	if encoding == "" {
		encoding = "utf-8"
	}
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=%s", mime, encoding))
	w.WriteHeader(code)
	w.Write(payload)
}
