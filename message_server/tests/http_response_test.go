package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"wraith.me/message_server/util"
)

func TestHttpErrRes(t *testing.T) {
	resp := util.ErrResponse(
		0,
		fmt.Errorf("error 1"),
		fmt.Errorf("error 2"),
		fmt.Errorf("error 3"),
	)
	fmt.Printf("%s\n", resp.MustJSON())
}

func TestHttpPayloadResp(t *testing.T) {
	resp := util.PayloadResponse(200, "many foos",
		foo1,
		foo2,
		foo3,
	)
	fmt.Printf("%s\n", resp.MustJSON())
}

// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go
func TestHttpResponder(t *testing.T) {
	//Setup handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		util.PayloadResponse(0, "this is a test", foo1, foo2, foo3).Respond(w)
	}

	//Setup a timeout
	//https://stackoverflow.com/a/55561566
	timeout := time.After(30 * time.Second)
	done := make(chan bool)
	go func() {
		http.HandleFunc("/", handler)
		err := http.ListenAndServe(":3333", nil)
		if err != nil {
			panic(err)
		}
	}()

	//Setup what to do once done
	select {
	case <-timeout:
	case <-done:
	}
}
