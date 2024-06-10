package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	ccrypto "wraith.me/message_server/crypto"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	c "wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/util"
	"wraith.me/message_server/util/httpu"
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
	ID mongoutil.UUID `json:"id" mapstructure:"id"`

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
Defines the structure of a user record that's returned from the database
during the existing user check. This includes the user's ID, public key,
and the status flags of the user.
*/
type existingUserResult struct {
	//`existingUserResult` extends the abstract entity type.
	obj.Entity `bson:",inline"`

	//The user's flags. These mark items such as verification status, deletion, etc.
	Flags obj.UserFlags `json:"flags" bson:"flags"`
}

/*
Handles incoming requests made to `POST /users/login_req`. This is stage 1
of the login process.
*/
func RequestLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new stage 1 object plus database result
	loginReq := loginUser{}
	hit := existingUserResult{}

	//Run pre-flight checks
	if !preFlight(&loginReq, &hit, w, r) {
		return
	}

	//Create a public key challenge using the user's info
	loginTok := c.NewPKChallenge(
		env.ID,
		hit.ID,
		c.CPurposeLOGIN,
		time.Now().Add(10*time.Minute),
		hit.Pubkey,
	).Encrypt(env.SK)

	//Send the token to the user
	httpu.HttpOkAsJson(w, loginTok, http.StatusOK)
}

/*
Handles incoming requests made to `POST /users/login_verify`. This is stage
2 of the login process.
*/
func VerifyLoginUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new stage 2 object plus database result
	loginVReq := loginVerifyUser{}
	hit := existingUserResult{}

	//Run pre-flight checks
	if !preFlight(&loginVReq, &hit, w, r) {
		return
	}

	//Verify the signature against the token; this proves ownership of the private key
	ok := ccrypto.Verify(loginVReq.PK, []byte(loginVReq.Token), loginVReq.Signature)
	if !ok {
		httpu.HttpErrorAsJson(w, fmt.Errorf("verification failure for PK %s against provided message and signature", loginVReq.PK.Fingerprint()), http.StatusForbidden)
		return
	}

	//Decrypt and validate the public key challenge
	//After this point, the user is considered fully authenticated; a token may now be issued
	loginTok, err := c.DecryptStrict(
		loginVReq.Token,
		env.SK,
		env.ID,
		loginVReq.ID,
		loginVReq.PK,
	)
	if err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusForbidden)
		return
	}

	fmt.Printf("verif_pk: %+v\n", loginTok)
	resp := fmt.Sprintf("REQUEST S2: %+v", loginVReq)
	httpu.HttpOkAsJson(w, resp, http.StatusOK)
}

// Contains the common FoC that is to be ran before any login request.
func preFlight[T loginUser | loginVerifyUser](user *T, hit *existingUserResult, w http.ResponseWriter, r *http.Request) bool {
	//Get the request body and attempt to parse from JSON
	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		httpu.HttpErrorAsJson(w, err, _PF_PARSE_ERR)
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
		httpu.HttpMultipleErrorsAsJson(w, missingErrors, _PF_PARSE_ERR)
		return false
	}

	//Ensure the incoming data is valid
	if err := ensureCorrectIdAndPKFmt(reqId, reqPK, reqSig); err != nil {
		httpu.HttpErrorAsJson(w, err, _PF_PARSE_ERR)
		return false
	}

	//Unmarshal the mapped request body into a user object
	if err := util.MSTextUnmarshal(reqBody, user); err != nil {
		httpu.HttpErrorAsJson(w, err, _PF_PARSE_ERR)
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
	//TODO: impl caching here
	userCollection := mcl.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)
	tmp, err := ensureExistantUser(userCollection, lu, r.Context())
	if err != nil {
		httpu.HttpErrorAsJson(w, err, _PF_PARSE_ERR)
		return false
	}

	//Check if a valid user was returned
	if hit == nil {
		httpu.HttpErrorAsJson(w, fmt.Errorf("no record found"), _PF_NO_USER)
		return false
	}
	*hit = *tmp

	//Check the user's flags to ensure they can actually sign-in
	//Their email and public key must be verified
	//TODO: Move this to auth if possible
	errors := []error{}
	if !hit.Flags.EmailVerified {
		errors = append(errors, fmt.Errorf("unverified email"))
	}
	if !hit.Flags.PubkeyVerified {
		errors = append(errors, fmt.Errorf("unverified public key"))
	}
	if len(errors) > 0 {
		httpu.HttpMultipleErrorsAsJson(w, errors, _PF_UNAUTHORIZED)
		return false
	}

	//No errors, so return true
	return true
}

// Ensures that the user ID and public key are of the proper format
func ensureCorrectIdAndPKFmt(id string, pk string, sig string) error {
	//Try to parse the UUID first
	if validId := mongoutil.IsValidUUIDv7(id); !validId {
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
func ensureExistantUser(coll *mongo.Collection, user loginUser, ctx context.Context) (*existingUserResult, error) {
	//Construct a Mongo aggregation pipeline to run the request; avoids making multiple round-trips to the database
	//This aggregation was exported from MongoDB; do not edit if you don't know what you are doing!
	agg := bson.A{
		//Match any documents that have the same ID and public key
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id", Value: user.ID},
					{Key: "pubkey", Value: user.PK},
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
	var hits []existingUserResult
	err := mongoutil.AggregateT(&hits, coll, agg, ctx)
	if err != nil {
		return nil, err
	}

	//Check if there was a hit
	if len(hits) > 0 {
		return &hits[0], nil
	}

	//Return nil for both since there was no record found, but no errors otherwise
	return nil, nil
}
