package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tanqiangyes/govalidator"
	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/bson"
	"wraith.me/message_server/crypto"
	"wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/obj/ip_addr"
	schema "wraith.me/message_server/schema/json"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/template/reg_email"
	"wraith.me/message_server/util"
)

/*
Represents a user object that was passed in as JSON. This object omits
stuff like `last_login`, `uuid`, `flags`, etc. Attributes correspond to
those on the standard user object. The `pubkey` attribute is serialized
as a base64 string to save on transport size.
*/
type intermediateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pubkey   string `json:"pubkey"`
}

/*
Represents a user object that is passed to the client once registration
is completed.
*/
type postsignupUser struct {
	//The ID of the user.
	ID util.UUID `json:"id"`

	//The username of the user.
	Username string `json:"username"`

	//The email of the user, but redacted.
	RedactedEmail string `json:"redacted_email"`

	//The fingerprint of the submitted public key.
	PKFingerprint string `json:"pk_fingerprint"`
}

// Handles incoming requests made to `POST /users/register`.
func RegisterUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new intermediate user object
	iuser := intermediateUser{}

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
	valid, verrs := iuser.validate(false)
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
	if err := ensureNonexistantUser(uc, iuser, r.Context()); err != nil {
		util.ErrResponse(http.StatusBadRequest, err).Respond(w)
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
		ip_addr.HttpIP2NetIP(r.RemoteAddr),
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
func ensureNonexistantUser(coll *user.UserCollection, usr intermediateUser, ctx context.Context) error {
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
	hits := make([]util.UUID, 0)
	err := coll.Aggregate(ctx, agg).All(&hits)
	if err != nil {
		return err
	}

	//Check if there were any hits
	if len(hits) > 0 {
		return fmt.Errorf("one or more provided fields already map to an existing user in the database")
	}
	return nil
}

// Validates an `intermediateUser` object using `tanqiangyes/govalidator`.
func (iu intermediateUser) validate(strictEmail bool) (bool, []error) {
	//Create a slice to hold the collected errors
	errors := []error{}

	//Step 1: Check the validity of the username
	//Username should be 4-16 characters in length and only consist of alphanumeric characters and underscores
	//This function should never throw an error since the regexp is hard-coded
	validUname, _ := regexp.MatchString(`^([a-z0-9_]){4,16}$`, strings.ToLower(iu.Username))
	if !validUname {
		errors = append(errors, fmt.Errorf("username '%s' is invalid. It must be 4-16 characters in length and only consist of alphanumeric characters and underscores", strings.ToLower(iu.Username)))
	}

	//Pick the appropriate email validator
	//`strictEmail` also ensures the email maps to an existing domain name
	emailValidator := govalidator.IsEmail[string]
	if strictEmail {
		emailValidator = govalidator.IsExistingEmail[string]
	}

	//Step 2: Check the validity of the email
	validEmail := emailValidator(strings.ToLower(iu.Email))
	if !validEmail {
		errors = append(errors, fmt.Errorf("email '%s' is invalid; it must be of the form 'foo@example.com'", strings.ToLower(iu.Email)))
	}

	//Step 3: Check the validity of the base64'ed public key by attempting to convert to a byte array
	validPubkey := true
	_, err := crypto.ParsePubkey(iu.Pubkey)
	if err != nil {
		validPubkey = false
		errors = append(errors, err)
	}

	//Return the validity status and any errors that occurred
	return validUname && validEmail && validPubkey, errors
}

/*
Performs post-signup operations on the newly created user object, such
as persistence to the database and generation of challenges.
*/
func postSignup(w http.ResponseWriter, r *http.Request, usr *user.User) error {
	//Issue a PASETO challenge for confirming the user's email
	paseto := challenge.NewEmailChallenge(
		env.ID,
		usr.ID,
		challenge.CPurposeCONFIRM,
		time.Now().Add(24*time.Hour),
		usr.Email,
	).Encrypt(env.SK)

	//Compose and send a challenge email to the user
	emailer := reg_email.NewRegEmail(
		*usr,
		util.TZOffsetFromReq(r),
		paseto,
		*cfg,
	)
	if err := emailer.Send(); err != nil {
		return err
	}

	//Persist the user in the database
	_, err := uc.InsertOne(r.Context(), usr)
	if err != nil {
		return err
	}

	//Write the response back to the user
	psu := postsignupUser{
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
