package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	ccrypto "wraith.me/message_server/crypto"
	c "wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

const (
	//The code to emit when the pre-flight parsing section fails.
	_PF_PARSE_ERR = http.StatusBadRequest

	//The HTTP status code to emit when the pre-flight existing user check fails.
	_PF_NO_USER = http.StatusNotFound

	//The HTTP status code to emit when the pre-flight authorization check fails.
	_PF_UNAUTHORIZED = http.StatusForbidden
)

/*
Defines the structure of JSON form data sent in the 1st stage of a login
request. This contains the user's ID and public key, both of which must
match what's in the database.
*/
type loginUser struct {
	//The UUID of the user to login as.
	ID util.UUID `json:"id" mapstructure:"id"`

	//The public key of the user to login as.
	PK ccrypto.Pubkey `json:"pk" mapstructure:"pk"`
}

/*
Defines the structure of JSON form data sent in the 2nd stage of a login
request. This contains everything that the 1st stage form data contains,
along with the token that was issued and the digital signature of the
token that was signed by the private key of the user.
*/
type loginVerifyUser struct {
	//`loginVerifyUser` extends `loginUser` by adding the previously generated token and the client's signature.
	loginUser `mapstructure:",squash"`

	//The login token that the user was given.
	Token string `json:"token" mapstructure:"token"`

	//The signature of the input token, signed by the user's private key.
	Signature ccrypto.Signature `json:"signature" mapstructure:"signature"`
}

/*
Handles incoming requests made to `POST /api/auth/login_req`. This is stage 1
of the login process.
*/
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new stage 1 object plus database result
	loginReq := loginUser{}
	user := user.User{}

	//Run pre-flight checks
	if !preFlight(&loginReq, &user, w, r) {
		return
	}

	//Create a public key challenge using the user's info
	loginTok := c.NewPKChallenge(
		env.ID,
		user.ID,
		c.CPurposeLOGIN,
		time.Now().Add(10*time.Minute),
		user.Pubkey,
	).Encrypt(env.SK)

	//Send the token to the user
	util.PayloadOkResponse(
		"",
		loginTok,
	).Respond(w)
}

/*
Handles incoming requests made to `POST /api/auth/login_verify`. This is stage
2 of the login process.
*/
func VerifyLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new stage 2 object plus database result
	loginVReq := loginVerifyUser{}
	user := user.User{}

	//Run pre-flight checks
	if !preFlight(&loginVReq, &user, w, r) {
		return
	}

	//Verify the signature against the token; this proves ownership of the private key
	ok := ccrypto.Verify(loginVReq.PK, []byte(loginVReq.Token), loginVReq.Signature)
	if !ok {
		util.ErrResponse(
			http.StatusForbidden,
			fmt.Errorf("verification failure for PK %s against provided message and signature", loginVReq.PK.Fingerprint()),
		).Respond(w)
		return
	}

	//Decrypt and validate the public key challenge
	//After this point, the user is considered fully authenticated; a token may now be issued
	loginTok, err := c.DecryptPKStrict(
		loginVReq.Token,
		env.SK,
		env.ID,
		c.CPurposeLOGIN,
		loginVReq.ID,
		loginVReq.PK,
	)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return
	}

	//TODO: mark user as verified and issue a login token here

	fmt.Printf("verif_pk: %+v\n", loginTok)
	util.PayloadOkResponse("", "ok").Respond(w)
}

// Contains the common FoC that is to be ran before any login request.
func preFlight[T loginUser | loginVerifyUser](user *T, hit *user.User, w http.ResponseWriter, r *http.Request) bool {
	//Get the request body and attempt to parse from JSON
	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		util.ErrResponse(_PF_PARSE_ERR, err).Respond(w)
		return false
	}

	//Ensure all request fields are present and are the correct type
	missingErrors := []error{}
	reqId, ok := reqBody["id"].(string)
	if !ok {
		missingErrors = append(missingErrors, fmt.Errorf("missing `id` field or it's poorly formed; expecting `string`"))
	}
	reqPK, ok := reqBody["pk"].(string)
	if !ok {
		missingErrors = append(missingErrors, fmt.Errorf("missing `pk` field or it's poorly formed; expecting `string`"))
	}

	//Check for the token and signature if this is a verification request
	var reqSig string
	if _, ok := any(user).(*loginVerifyUser); ok {
		if _, ok := reqBody["token"].(string); !ok {
			missingErrors = append(missingErrors, fmt.Errorf("missing `token` field or it's poorly formed; expecting `string`"))
		}
		reqSig, ok = reqBody["signature"].(string)
		if !ok {
			missingErrors = append(missingErrors, fmt.Errorf("missing `signature` field or it's poorly formed; expecting `string`"))
		}
	}

	//Error out if any required field is missing
	if len(missingErrors) > 0 {
		util.ErrResponse(_PF_PARSE_ERR, missingErrors...).Respond(w)
		return false
	}

	//Ensure the incoming data is valid
	if err := ensureCorrectIdAndPKFmt(reqId, reqPK, reqSig); err != nil {
		util.ErrResponse(_PF_PARSE_ERR, err).Respond(w)
		return false
	}

	//Unmarshal the mapped request body into a user object
	if err := util.MSTextUnmarshal(reqBody, user, ""); err != nil {
		util.ErrResponse(_PF_PARSE_ERR, err).Respond(w)
		return false
	}

	//Derive a common "loginUser" type; Go generics are kinda dumb
	//See: https://go.dev/play/p/H3fBSekLyE6
	var lu loginUser
	switch any(user).(type) {
	case *loginUser:
		tmp, _ := any(user).(*loginUser)
		lu = *tmp
	case *loginVerifyUser:
		tmp, _ := any(user).(*loginVerifyUser)
		lu = tmp.loginUser
	default:
		panic(fmt.Sprintf("Unexpected type %T\n", user)) //This block shouldn't ever be hit
	}

	//Ensure the claims map to an existing user in the database
	tmp, err := ensureExistantUser(uc, lu, r.Context())
	if err != nil {
		//Check if the error has to do with a missing user
		code := _PF_PARSE_ERR
		desc := err
		if errors.Is(err, mongo.ErrNoDocuments) {
			code = _PF_NO_USER
			desc = fmt.Errorf("no such user with ID %s", lu.ID)
		}

		//Respond back with the error
		util.ErrResponse(code, desc).Respond(w)
		return false
	}

	//Check if a valid user was returned
	/*
		if tmp == nil {
			util.ErrResponse(
				_PF_NO_USER,
				fmt.Errorf("no such user with ID %s", lu.ID),
			).Respond(w)
			return false
		}
	*/
	*hit = *tmp

	//Check the user's flags to ensure they can actually sign-in
	//Their email and public key must be verified
	//TODO: Move this to auth if possible
	//TODO: login verifies public key automatically
	errors := []error{}
	if !hit.Flags.EmailVerified {
		errors = append(errors, fmt.Errorf("unverified email"))
	}
	if !hit.Flags.PubkeyVerified {
		errors = append(errors, fmt.Errorf("unverified public key"))
	}
	if len(errors) > 0 {
		util.ErrResponse(_PF_UNAUTHORIZED, errors...).Respond(w)
		return false
	}

	//No errors, so return true
	return true
}

// Ensures that the user ID and public key are of the proper format
func ensureCorrectIdAndPKFmt(id string, pk string, sig string) error {
	//Try to parse the UUID first
	if validId := util.IsValidUUIDv7(id); !validId {
		return fmt.Errorf("invalid UUID format `%s`; expected a UUIDv7 in the form: `xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx`", id)
	}

	//Check the validity of the base64'ed public key by attempting to convert to a byte array
	dbytes, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return fmt.Errorf("request.pk: %s", err)
	}
	if len(dbytes) != ccrypto.PUBKEY_SIZE {
		return fmt.Errorf("mismatched public key size (%d); expected: %d", len(dbytes), ccrypto.PUBKEY_SIZE)
	}

	//Check the validity of the base64'ed signature by attempting to convert to a byte array
	//Only do this if its not empty
	if len(sig) > 0 {
		sbytes, err := base64.StdEncoding.DecodeString(sig)
		if err != nil {
			return fmt.Errorf("request.signature: %s", err)
		}
		if len(sbytes) != ccrypto.SIG_SIZE {
			return fmt.Errorf("mismatched signature size (%d); expected: %d", len(sbytes), ccrypto.SIG_SIZE)
		}
	}

	//No errors so return nil
	return nil
}

// Ensures that a user with the given UUID and public key exists.
func ensureExistantUser(coll *user.UserCollection, usr loginUser, ctx context.Context) (*user.User, error) {
	//Construct a Mongo aggregation pipeline to run the request; avoids making multiple round-trips to the database
	//This aggregation was exported from MongoDB; do not edit if you don't know what you are doing!
	agg := bson.A{
		//Match any documents that have the same ID and public key
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id", Value: usr.ID},
					{Key: "pubkey", Value: usr.PK},
				},
			},
		},
	}

	//Run the request and collect all hits; critical errors may be reported from this function so handle appropriately
	var hit user.User
	err := coll.Aggregate(ctx, agg).One(&hit)
	if err != nil {
		return nil, err
	}

	//Check if there was a hit
	if hit.ID != util.NilUUID() {
		return &hit, nil
	}

	//Return nil for both since there was no record found, but no errors otherwise
	return nil, nil
}
