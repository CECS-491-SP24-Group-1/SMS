package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/pkg/controller/csolver"
	"wraith.me/message_server/pkg/crypto"
	"wraith.me/message_server/pkg/db/mongoutil"
	"wraith.me/message_server/pkg/http_types/request"
	"wraith.me/message_server/pkg/http_types/response"
	"wraith.me/message_server/pkg/obj/ip_addr"
	schema "wraith.me/message_server/pkg/schema/json"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

// Handles incoming requests made to `POST /api/auth/register`.
func RegisterUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new intermediate user object
	iuser := request.RegisteringUser{}

	//Read in the request body to a string
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Validate the request body against the registration JSON schema
	result, err := gojsonschema.Validate(schema.Register, gojsonschema.NewBytesLoader(bodyBytes))
	if err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w)
		return
	}
	if !result.Valid() {
		//Collect the validation errors and report them to the client
		verrs := make([]error, len(result.Errors()))
		for i, err := range result.Errors() {
			verrs[i] = fmt.Errorf(err.Description())
		}
		util.ErrResponse(http.StatusBadRequest, verrs...).Respond(w)
		return
	}

	//Get the request body and attempt to parse to JSON
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&iuser); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w)
		return
	}

	//Validate the unmarshalled struct
	//At this point, the incoming JSON was accepted, but fields may be missing or invalid
	valid, verrs := iuser.Validate(false)
	if !valid {
		util.ErrResponse(http.StatusBadRequest, verrs...).Respond(w)
		return
	}

	//Decode the base64 public key to a byte array
	decodedPK, err := base64.StdEncoding.DecodeString(iuser.Pubkey)
	if err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w)
		return
	}

	//Ensure the user doesn't already exist in the database
	exists, err := ensureNonexistantUser(uc, iuser, r.Context())
	if err != nil {
		fmt.Printf("error during request from %s: %s\n", r.Host, err)
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}

	//Check if there were any hits
	if exists {
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("one or more provided fields already map to an existing user in the database"),
		).Respond(w)
		return
	}

	//Fill in the rest of the details
	uuid := util.MustNewUUID7()
	user := user.NewUser(
		uuid,
		crypto.NilPubkey(),
		strings.ToLower(iuser.Username),
		iuser.Username,
		strings.ToLower(iuser.Email),
		util.NowMillis(),
		ip_addr.HttpIP2IPAddr(r.RemoteAddr),
		user.DefaultUserFlags(),
		user.DefaultUserOptions(),
	)
	copy(user.Pubkey[:], decodedPK[:])

	//Complete the post-signup steps, including challenge generation and issuance of a temporary token
	if err := postSignup(w, r, user); err != nil {
		util.ErrResponse(http.StatusInternalServerError, err).Respond(w)
		return
	}
}

/*
Ensures that a user doesn't already exist in the database based on what
was given by the user. A `nil` error indicates that no matching records
were found. Checking collections for existent objects is expensive, so
not all records are checked if one fails.
*/
func ensureNonexistantUser(coll *user.UserCollection, usr request.RegisteringUser, ctx context.Context) (bool, error) {
	//Parse out the public key of the incoming user
	pubkey, _ := crypto.ParsePubkey(usr.Pubkey) //Errors should not occur here; data is already pre-validated

	//Construct a Mongo aggregation pipeline to run the request; avoids making multiple round-trips to the database
	//This aggregation was exported from MongoDB; do not edit if you don't know what you are doing!
	agg := bson.A{
		//Match any documents that have the same username, email, or public key
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "$or",
						Value: bson.A{
							bson.D{{Key: "username", Value: usr.Username}},
							bson.D{{Key: "email", Value: usr.Email}},
							bson.D{{Key: "pubkey", Value: pubkey}},
						},
					},
				},
			},
		},
		//Reduce the size of the incoming BSON documents to improve performance
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: 1},
				},
			},
		},
	}

	//Run the request and collect all hits; critical errors may be reported from this function so handle appropriately
	hits, err := mongoutil.Aggregate2IDArr(coll.Aggregate(ctx, agg))
	if err != nil {
		return false, err
	}

	return len(hits) > 0, nil
}

/*
Performs post-signup operations on the newly created user object, such
as persistence to the database and generation of challenges.
*/
func postSignup(w http.ResponseWriter, r *http.Request, usr *user.User) error {
	//Issue an email challenge for the user if email is enabled
	if cfg.Email.Enabled {
		if err := csolver.IssueEmailChallenge(usr, cfg, env, r); err != nil {
			return err
		}
	} else {
		//Email is not enabled, so their email is marked verified by default
		usr.Flags.EmailVerified = true
	}

	//Persist the user in the database
	_, err := uc.InsertOne(r.Context(), usr)
	if err != nil {
		return err
	}

	//Write the response back to the user
	psu := response.RegisteredUser{
		ID:            usr.ID,
		Username:      usr.Username,
		RedactedEmail: util.RedactEmail(usr.Email),
		PKFingerprint: usr.Pubkey.Fingerprint(),
	}
	util.PayloadResponse(
		http.StatusCreated,
		fmt.Sprintf("created new user with ID %s", psu.ID),
		psu,
	).Respond(w)

	//No errors so return nil
	return nil
}
