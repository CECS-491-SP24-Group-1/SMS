package obj

//
//-- CLASS: Response
//

/*
Represents a status message that's returned to a client by the server.
This is usually done if the server wants to send feedback along with
the HTTP status code, such as when an error occurs for whatever reason.
*/
type Response struct {
	//The HTTP status code of the response.
	Code int `json:"code"`

	//The status indicator that is to be sent with the message.
	Status string `json:"status"`

	//The description of the status that is to be sent.
	Desc string `json:"desc"`
}

//
//-- CLASS: MultiResponse
//

/*
Represents a status message that's returned to a client by the server.
This response is like a regular `Response`, but allows for multiple
messages to be sent. This may be desirable if, for example, multiple
errors are to be returned to the client.
*/
type MultiResponse struct {
	//The HTTP status code of the response.
	Code int `json:"code"`

	//The status indicator that is to be sent with the message.
	Status string `json:"status"`

	//The description messages of the status that is to be sent.
	Descs []string `json:"descs"`
}

// Represents an object that's returned to a client when an error occurs.
// TODO: deprecate in favor of `Response`.
type ErrorResp struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Error  string `json:"error"`
}

// Represents an object that's returned to a client when multiple errors occur.
// TODO: deprecate in favor of `MultiResponse`.
type MultiErrorResp struct {
	Status string   `json:"status"`
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}
