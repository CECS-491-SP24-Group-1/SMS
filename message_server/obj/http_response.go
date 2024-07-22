package obj

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//
//-- CLASS: HttpResponse
//

// Represents a response that is sent back to the user after an API call.
type HttpResponse[T any] struct {
	//The HTTP code to emit; default `200`.
	Code int `json:"code"`

	//The status message to emit; default: <http response description>.
	Status string `json:"status,omitempty"`

	//The description of the response; default: `""`.
	Desc string `json:"desc,omitempty"`

	//The errors of the server response, if any; default: `[]`.
	Errors []string `json:"errors,omitempty"` //https://github.com/golang/go/issues/5161

	//The payloads of the server response, if any; default: `nil`.
	Payload []T `json:"payload,omitempty"`
}

// -- Constructors
// Creates a new error response from a list of errors.
func ErrResponse(code int, errs ...error) HttpResponse[bool] {
	//Determine the correct error code
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	//Create the status message
	estr := "error"
	if len(errs) > 1 {
		estr += "s"
	}
	status := fmt.Sprintf("%s; %d %s", http.StatusText(code), len(errs), estr)

	//Create the response
	resp := HttpResponse[bool]{
		Code:   code,
		Status: status,
		Errors: make([]string, len(errs)),
	}

	//Unpack the errors
	for i, err := range errs {
		resp.Errors[i] = err.Error()
	}

	//Emit the full response
	return resp
}

// Creates a new info response from an HTTP code and description.
func InfoResponse(code int, desc string) HttpResponse[bool] {
	//Determine the correct response code
	if code <= 0 {
		code = http.StatusOK
	}

	//Construct and return the response object
	return HttpResponse[bool]{
		Code:   code,
		Status: "ok",
		Desc:   desc,
	}
}

// Creates a new ok response from a description.
func OkResponse(desc string) HttpResponse[bool] {
	return InfoResponse(http.StatusOK, desc)
}

// Creates a new payload response from a single payload object or multiple.
func PayloadResponse[T any](code int, desc string, payload ...T) HttpResponse[T] {
	//Determine the correct response code
	if code <= 0 {
		code = http.StatusOK
	}

	//Create the status message
	pstr := "payload"
	if len(payload) > 1 {
		pstr += "s"
	}
	status := fmt.Sprintf("%s; %d %s", http.StatusText(code), len(payload), pstr)

	//Create the response
	resp := HttpResponse[T]{
		Code:    code,
		Status:  status,
		Desc:    desc,
		Payload: payload,
	}

	//Emit the full response
	return resp
}

// Creates a new error payload response from a single payload object or multiple.
func PayloadErrResponse[T any](desc string, payload ...T) HttpResponse[T] {
	return PayloadResponse[T](http.StatusInternalServerError, desc, payload...)
}

// Creates a new ok payload response from a single payload object or multiple.
func PayloadOkResponse[T any](desc string, payload ...T) HttpResponse[T] {
	return PayloadResponse[T](http.StatusOK, desc, payload...)
}

// -- Methods
// Emits the JSON encoding of this object, along with any errors that occurred.
func (r HttpResponse[T]) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// Emits the JSON encoding of this object, panicking if any errors occur.
func (r HttpResponse[T]) MustJSON() []byte {
	json, err := r.JSON()
	if err != nil {
		panic(fmt.Sprintf("HttpResponse::MustJson(); %s", err))
	}
	return json
}

// Writes the JSON encoding of this object to an HTTP response writer.
func (r HttpResponse[T]) Respond(w http.ResponseWriter) {
	//Set initial headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	//Marshal to JSON
	json, err := r.JSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"code": %d, "errors": ["failed to marshal HttpResponse to json; %s"]}`, http.StatusInternalServerError, err)
		return
	}

	//Write the status code and resultant JSON to the output stream
	w.WriteHeader(r.Code)
	_, err = w.Write(json)
	if err != nil {
		panic(fmt.Sprintf("HttpResponse::Respond(); %s", err))
	}
}
