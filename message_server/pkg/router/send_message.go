package router

import (
	"fmt"
	"io"
	"net/http"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("u got mail::%s\n", string(body))
}
