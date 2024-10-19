package tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	ccrypto "wraith.me/message_server/pkg/crypto"
	"wraith.me/message_server/pkg/obj/token"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
	"wraith.me/message_server/pkg/util/ms"
)

var usr = user.NewUserSimple(
	"johndoe123",
	"johndoeismyname@example.com",
)

func TestUser2BSON(t *testing.T) {
	//Create a mock UUID and key for the issuer
	issuerId := util.MustNewUUID7()
	_, issuerKey, err := ccrypto.NewKeypair(nil)
	if err != nil {
		t.Fatal(err)
	}

	//Add a test token
	//Decryption error will occur since the encryption key is different than prod
	tok := token.NewToken(
		usr.ID,
		issuerId,
		token.TokenTypeREFRESH,
		time.Now().Add(5*time.Minute),
		nil,
		nil,
	)
	tstr := tok.Encrypt(issuerKey, true)
	usr.AddToken(tok.ID.String(), tstr, tok.Expiry)

	//Marshal to BSON
	bb, err := bson.Marshal(usr)
	if err != nil {
		t.Fatal(err)
	}

	//Unmarshal the BSON back to an object
	var out user.User
	if err := bson.Unmarshal(bb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *usr)
	fmt.Printf("OUT: %v\n", out)
}

func TestUser2JSON(t *testing.T) {
	//Create a mock UUID and key for the issuer
	issuerId := util.MustNewUUID7()
	_, issuerKey, err := ccrypto.NewKeypair(nil)
	if err != nil {
		t.Fatal(err)
	}

	//Add a test token
	//Decryption error will occur since the encryption key is different than prod
	tok := token.NewToken(
		usr.ID,
		issuerId,
		token.TokenTypeREFRESH,
		time.Now().Add(5*time.Minute),
		nil,
		nil,
	)
	tstr := tok.Encrypt(issuerKey, true)
	usr.AddToken(tok.ID.String(), tstr, tok.Expiry)

	//Marshal to JSON
	jb, err := json.Marshal(usr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("JSON: %s\n", jb) //Tokens are not marshaled with the JSON

	//Unmarshal the JSON back to an object
	var out user.User
	if err := json.Unmarshal(jb, &out); err != nil {
		t.Fatal(err)
	}

	//Print the input and output objects
	fmt.Printf("IN:  %v\n", *usr)
	fmt.Printf("OUT: %v\n", out)
}

func TestUser2MS(t *testing.T) {
	//Create a mock user
	user := user.NewUserSimple("johndoe", "johndoe@example.com")
	fmt.Printf("struct: %v\n", user)

	//Marshal to a map using mapstructure
	mms := make(map[string]interface{})
	if err := ms.MSTextMarshal(*user, &mms, "bson"); err != nil {
		t.Fatal(err)
	}

	//Redact some fields as a test
	mms["id"] = mms["UUID"]
	delete(mms, "UUID")
	delete(mms, "flags")
	delete(mms, "last_ip")
	delete(mms, "options")
	delete(mms, "tokens")

	//Emit the map to stdout
	fmt.Printf("mapstr: %v\n", mms)
}
