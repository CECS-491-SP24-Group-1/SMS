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

	//The description message of the status that is to be sent.
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

	//The description messages of the statuses that are to be sent.
	Desc []string `json:"desc"`
}
