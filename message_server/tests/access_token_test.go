package tests

import (
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/obj/token"
)

func TestGenAccessTok(t *testing.T) {
	//Get the server secrets
	secrets, err := config.EnvInit("../secrets.env")
	if err != nil {
		t.Fatal(err)
	}

	//Get a random user
	usr, err := GetRandomUser()
	if err != nil {
		t.Fatal(err)
	}

	//Generate an access token
	//The tokens generated should work for auth tests, but are not persisted to the db
	tok := token.NewToken(
		usr.ID,
		secrets.ID,
		token.TokenTypeACCESS,
		time.Now().Add(30*time.Minute),
		nil,
		nil,
	).Encrypt(secrets.SK, true)

	//Get the expiration
	expr, err := token.GetExprFromFooter(tok)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("User UUID:  %s\n", usr.ID)
	fmt.Printf("User Token: %s\n", tok)
	fmt.Printf("Token Expr: %s\n", expr.Format(token.TimeFmt))
}
