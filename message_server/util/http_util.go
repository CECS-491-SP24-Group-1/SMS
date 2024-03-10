package util

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"wraith.me/message_server/obj"
)

// Writes plaintext to an HTTP response and issues a `500` status code.
func HttpErrorSimple(w http.ResponseWriter, msg string) {
	TextualResponse(w, []byte(msg), http.StatusInternalServerError, "", "")
}

// Writes JSON to a request and issues a `500` status code.
func HttpErrorJson(w http.ResponseWriter, json string) {
	TextualResponse(w, []byte(json), http.StatusInternalServerError, "application/json", "")
}

// Writes plaintext to an HTTP response and issues a `200` status code.
func HttpOkSimple(w http.ResponseWriter, msg string) {
	TextualResponse(w, []byte(msg), http.StatusOK, "", "")
}

// Writes JSON to a request and issues a `200` status code.
func HttpOkJson(w http.ResponseWriter, json string) {
	TextualResponse(w, []byte(json), http.StatusOK, "application/json", "")
}

// Issues an HTTP error response as JSON.
func HttpErrorAsJson(w http.ResponseWriter, err error, code int) {
	//Construct an error response object
	resp := obj.ErrorResp{
		Code:   code,
		Status: "error",
		Error:  err.Error(),
	}

	//Write the object to the response writer
	w.WriteHeader(code)
	jerr := json.NewEncoder(w).Encode(&resp)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
	}
}

// Issues an HTTP error response as JSON.
func HttpMultipleErrorsAsJson(w http.ResponseWriter, errs []error, code int) {
	//Populate a list of stringified errors
	stringErrs := make([]string, len(errs))
	for i, err := range errs {
		stringErrs[i] = err.Error()
	}

	//Construct an error response object
	resp := obj.MultiErrorResp{
		Code:   code,
		Status: "multiple_errors",
		Errors: stringErrs,
	}

	//Write the object to the response writer
	w.WriteHeader(code)
	jerr := json.NewEncoder(w).Encode(&resp)
	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
	}
}

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

func HttpIP2NetIP(ip string) net.IP {
	rawIP := ip[0:strings.LastIndex(ip, ":")] //Get just the IP; last colon indicates the port
	return net.ParseIP(rawIP)
}
