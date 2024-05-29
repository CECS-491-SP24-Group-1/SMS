package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/util/httpu"
)

// Sets the key name for the subject public key field in a PASETO token.
const SUBJECT_PK_KEY = "sub-pk"

/*
Defines the structure of JSON form data sent in the 1st stage of a login
request. This contains the user's ID and public key, both of which must
match what's in the database.
*/
type loginUser struct {
	//The UUID of the user to login as.
	ID string `json:"id"`

	//THe public key of the user to login as.
	PK string `json:"pk"`
}

/*
Defines the structure of JSON form data sent in the 2nd stage of a login
request. This contains everything that the 1st stage form data contains,
along with the token that was issued and the digital signature of the
token that was signed by the private key of the user.
*/
type loginVerifyUser struct {
	//`loginVerifyUser` extends `loginUser` by adding the previously generated token and the client's signature.
	loginUser

	//The login token that the user was given.
	Token string `json:"token"`

	//The signature of the input token, signed by the user's private key.
	Signature string `json:"signature"`
}

type existingUserResult struct {
}

//TODO: create pre-login helper function, which contains the common FoC between the request and verify routes

// Handles incoming requests made to `POST /users/login_req`.
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new stage 1 object
	loginReq := loginUser{}

	//Get the request body and attempt to parse to JSON
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Ensure the incoming data is valid
	if err := ensureCorrectIdAndPKFmt(loginReq); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Ensure the claims map to an existing user in the database

	resp := fmt.Sprintf("REQUEST S1: %+v", loginReq)
	w.Write([]byte(resp))
}

// Handles incoming requests made to `POST /users/login_verify`.
func VerifyLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	// Create a new stage 2 object
	loginVReq := loginVerifyUser{}

	//Get the request body and attempt to parse to JSON
	if err := json.NewDecoder(r.Body).Decode(&loginVReq); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	fmt.Printf("verif_pk: %+v\n", loginVReq.loginUser)
	resp := fmt.Sprintf("REQUEST S2: %+v", loginVReq)
	w.Write([]byte(resp))
}

// Ensures that a user with the given UUID and public key exists.
func ensureExistantUser(coll *mongo.Collection, user loginUser, ctx context.Context) (bool, error) {
	//Parse out the ID public key of the incoming user
	id := mongoutil.UUIDFromString(user.ID)
	pubkey, _ := obj.ParsePubkeyBytes(user.PK) //Errors should not occur here; data is already pre-validated

	//Construct a Mongo aggregation pipeline to run the request; avoids making multiple round-trips to the database
	//This aggregation was exported from MongoDB; do not edit if you don't know what you are doing!
	agg := bson.A{
		//Match any documents that have the same ID and public key
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id", Value: id},
					{Key: "pubkey", Value: pubkey},
				},
			},
		},
		//Reduce the size of the incoming BSON documents to improve performance, but leave the flags intact
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 1},
					{Key: "pubkey", Value: 1},
					{Key: "flags", Value: 1},
				},
			},
		},
	}

	//Run the request and collect all hits; critical errors may be reported from this function so handle appropriately
	hits, err := mongoutil.Aggregate(coll, agg, ctx)
	if err != nil {
		return false, err
	}

	//Check if there were any hits
	return len(hits) > 0, nil
}

// Ensures that the user ID and public key are of the proper format
func ensureCorrectIdAndPKFmt(user loginUser) error {
	//Try to parse the UUID first
	if validId := mongoutil.IsValidUUIDv7(user.ID); !validId {
		return fmt.Errorf("invalid UUID format `%s`; expected a UUIDv7 in the form: `xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`", user.ID)
	}

	//Check the validity of the base64'ed public key by attempting to convert to a byte array
	dbytes, err := base64.StdEncoding.DecodeString(user.PK)
	if err != nil {
		return err
	}
	validPubkey := len(dbytes) == obj.PUBKEY_SIZE
	if !validPubkey {
		return fmt.Errorf("mismatched public key size (%d); expected: %d", len(dbytes), obj.PUBKEY_SIZE)
	}

	//No errors so return nil
	return nil
}
