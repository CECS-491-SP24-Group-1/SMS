package csolver

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/config"
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
)

// Issues a public key challenge for a user. This is stage 1 of a login/pk challenge.
func IssuePKChallenge(user user.User, env *config.Env) string {
	return c.NewPKChallenge(
		env.ID,
		user.ID,
		c.CPurposeLOGIN,
		time.Now().Add(10*time.Minute),
		user.Pubkey,
	).Encrypt(env.SK)
}

// Verifies that a public key challenge is valid. This is stage 2 of a login/pk challenge.
func VerifyPKChallenge(vreq LoginVerifyUser, env *config.Env) (*c.CToken, error) {
	//Verify the signature against the token; this proves ownership of the private key
	ok := ccrypto.Verify(vreq.PK, []byte(vreq.Token), vreq.Signature)
	if !ok {
		return nil, fmt.Errorf("verification failure for PK %s against provided message and signature", vreq.PK.Fingerprint())
	}

	//Decrypt and validate the public key challenge
	//After this point, the user is considered fully authenticated; a token may now be issued
	loginTok, err := c.DecryptPKStrict(
		vreq.Token,
		env.SK,
		env.ID,
		c.CPurposeLOGIN,
		vreq.ID,
		vreq.PK,
	)
	if err != nil {
		return nil, err
	}

	//Double check to ensure the challenge PK and the user PK match up
	if subtle.ConstantTimeCompare(
		[]byte(vreq.PK.String()), []byte(loginTok.Claim),
	) == 0 {
		return nil, fmt.Errorf("token claim and user key mismatch")
	}

	//Return the login token
	return loginTok, nil
}

// Contains the common FoC that is to be ran before any login/pubkey solve request.
func PreFlight[T LoginUser | LoginVerifyUser](user *T, hit *user.User, uc *user.UserCollection, w http.ResponseWriter, r *http.Request) bool {
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
	if _, ok := any(user).(*LoginVerifyUser); ok {
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
	var lu LoginUser
	switch any(user).(type) {
	case *LoginUser:
		tmp, _ := any(user).(*LoginUser)
		lu = *tmp
	case *LoginVerifyUser:
		tmp, _ := any(user).(*LoginVerifyUser)
		lu = tmp.LoginUser
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
func ensureExistantUser(coll *user.UserCollection, usr LoginUser, ctx context.Context) (*user.User, error) {
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
