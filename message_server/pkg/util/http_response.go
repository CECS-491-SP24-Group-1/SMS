package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	//Controls whether singular errors and payloads should be marshalled as arrays (default: true).
	MarshalSingularAsArrays = true
)

//
//-- CLASS: HttpResponse
//

// Represents a response that is sent back to the user after an API call.
type HttpResponse[T any] struct {
	//The HTTP code to emit; default `200` for info/ok/payload, `500` for error.
	Code int `json:"code"`

	//The status message to emit; default: `<http response description>`.
	Status string `json:"status"`

	//The description of the response; default: `""`.
	Desc string `json:"desc"`

	//The errors of the server response, if any; default: `[]`.
	Errors []string `json:"errors,omitempty"` //https://github.com/golang/go/issues/5161

	//The payloads of the server response, if any; default: `[]`.
	Payloads []T `json:"payloads,omitempty"`

	//Whether the response has a payload.
	hasPayload bool
}

//-- Constructors

// Creates a new error response from a list of errors.
func ErrResponse(code int, errs ...error) HttpResponse[bool] {
	//Determine the correct error code
	if code <= 0 {
		code = http.StatusInternalServerError
	}

	//Create the status message
	desc := fmt.Sprintf("%d error", len(errs))
	desc += If(len(errs) > 1, "s", "")

	//Create the response
	resp := HttpResponse[bool]{
		Code:   code,
		Status: http.StatusText(code) + "; " + desc,
		Desc:   desc,
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
		Status: http.StatusText(code),
		Desc:   desc,
	}
}

// Creates a new ok response from a description.
func OkResponse(desc string) HttpResponse[bool] {
	return InfoResponse(http.StatusOK, desc)
}

// Creates a new payload response from a single payload object or multiple.
func PayloadResponse[T any](code int, desc string, payloads ...T) HttpResponse[T] {
	//Determine the correct response code
	if code <= 0 {
		code = http.StatusOK
	}

	//Create the status message
	pstr := "payload"
	if len(payloads) > 1 {
		pstr += "s"
	}
	status := fmt.Sprintf("%s; %d %s", http.StatusText(code), len(payloads), pstr)

	//Add the description if it's missing
	if strings.TrimSpace(desc) == "" {
		desc = fmt.Sprintf("%T; x%d", payloads, len(payloads))
	}

	//Create the response
	resp := HttpResponse[T]{
		Code:       code,
		Status:     status,
		Desc:       desc,
		Payloads:   payloads,
		hasPayload: true,
	}

	//Emit the full response
	return resp
}

// Creates a new error payload response from a single payload object or multiple.
func PayloadErrResponse[T any](desc string, payloads ...T) HttpResponse[T] {
	return PayloadResponse[T](http.StatusInternalServerError, desc, payloads...)
}

// Creates a new ok payload response from a single payload object or multiple.
func PayloadOkResponse[T any](desc string, payloads ...T) HttpResponse[T] {
	return PayloadResponse[T](http.StatusOK, desc, payloads...)
}

//-- Methods

// Emits the JSON encoding of this object, along with any errors that occurred.
func (r HttpResponse[T]) JSON() ([]byte, error) {
	//Get the inner type of the payload and don't escape it if it's a string
	stringPayloads, ok := any(r.Payloads).([]string)
	if ok {
		//Marshal the strings to raw JSON
		jsons := make([]json.RawMessage, len(r.Payloads))
		for i := range jsons {
			//Get the current payload
			pload := stringPayloads[i]

			//If the payload is valid JSON, encode it normally, otherwise return it as a regular JSON string
			if json.Valid([]byte(pload)) {
				jsons[i] = json.RawMessage(pload)
			} else {
				jsons[i] = json.RawMessage(fmt.Sprintf("%q", pload))
			}
		}

		//Create a new object from the existing one
		obj := HttpResponse[json.RawMessage]{
			Code:     r.Code,
			Status:   r.Status,
			Desc:     r.Desc,
			Errors:   r.Errors,
			Payloads: jsons,
		}

		//Marshal the new object as usual
		return backendMarshal(&obj)
	}

	//Marshal normally for every other type
	return backendMarshal(&r)
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
		fmt.Fprintf(w,
			`{"code": %d, "status": "error during JSON marshal", "errors": ["failed to marshal HttpResponse to json; %s"]}`,
			http.StatusInternalServerError, err)
		return
	}

	//fmt.Printf("output JSON: %s\n", json)

	//Write the status code and resultant JSON to the output stream
	w.WriteHeader(r.Code)
	_, err = w.Write(json)
	if err != nil {
		panic(fmt.Sprintf("HttpResponse::Respond(); %s", err))
	}
}

//-- Private utilities

// Handles the backend of the JSON marshalling operation.
func backendMarshal[T any](obj *HttpResponse[T]) ([]byte, error) {
	//Create an alias of the object
	type alias HttpResponse[T]
	aux := struct {
		*alias
		Error    interface{} `json:"error,omitempty"`
		Errors   interface{} `json:"errors,omitempty"`
		Payload  interface{} `json:"payload,omitempty"`
		Payloads interface{} `json:"payloads,omitempty"`
	}{
		alias: (*alias)(obj),
	}

	//Marshal singular items into an array or separately in singular fields
	if MarshalSingularAsArrays {
		//Add the errors and payloads only if they are non-empty
		if len(obj.Errors) > 0 {
			aux.Errors = obj.Errors
		}

		//Always include payloads if hasPayload is true
		if obj.hasPayload {
			aux.Payloads = obj.Payloads
		}
	} else {
		//Add the errors and payloads only if they are non-empty
		if len(obj.Errors) == 1 {
			aux.Error = obj.Errors[0]
		} else if len(obj.Errors) > 1 {
			aux.Errors = obj.Errors
		}

		//Only add payloads if hasPayload is true
		if obj.hasPayload {
			if len(obj.Payloads) == 1 {
				aux.Payload = obj.Payloads[0]
			} else {
				aux.Payloads = obj.Payloads
			}
		}
	}

	//Marshal the alias
	return json.Marshal(&aux)
}
