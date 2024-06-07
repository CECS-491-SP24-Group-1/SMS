package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/tanqiangyes/govalidator"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/crud"
	"wraith.me/message_server/crypto"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/email"
	"wraith.me/message_server/mw"
	"wraith.me/message_server/obj"
	"wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/obj/ip_addr"
	cr "wraith.me/message_server/redis"
	remailt "wraith.me/message_server/template/registration_email"
	"wraith.me/message_server/util"
	"wraith.me/message_server/util/httpu"
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
	ID mongoutil.UUID `json:"id"`

	//The email of the user, but redacted.
	RedactedEmail string `json:"redacted_email"`

	//The IDs of the challenges that the user must fulfil for registration to be completed.
	Challenges []mongoutil.UUID `json:"challenges"`

	//A token used to allow temporary API access to solve challenges. This key is only valid for that endpoint.
	TempAccessToken string `json:"temp_access_token"`
}

// Handles incoming requests made to `POST /users/register`.
func RegisterUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new intermediate user object
	iuser := intermediateUser{}

	//Get the request body and attempt to parse to JSON
	if err := json.NewDecoder(r.Body).Decode(&iuser); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Validate the unmarshalled struct
	//At this point, the incoming JSON was accepted, but fields may be missing or invalid
	valid, verrs := iuser.validate(false)
	if !valid {
		httpu.HttpMultipleErrorsAsJson(w, verrs, http.StatusBadRequest)
		return
	}

	//Decode the base64 public key to a byte array
	decodedPK, err := base64.StdEncoding.DecodeString(iuser.Pubkey)
	if err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Get the users collection from the database and ensure a record doesn't already exist
	userCollection := mcl.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)

	//Ensure the user doesn't already exist in the database
	if err := ensureNonexistantUser(userCollection, iuser, r.Context()); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Fill in the rest of the details
	uuid := mongoutil.MustNewUUID7()
	user := obj.NewUser(
		uuid,
		crypto.NilPubkey(),
		strings.ToLower(iuser.Username),
		iuser.Username,
		strings.ToLower(iuser.Email),
		util.NowMillis(),
		ip_addr.HttpIP2NetIP(r.RemoteAddr),
		obj.DefaultUserFlags(),
		obj.DefaultUserOptions(),
	)
	copy(user.Pubkey[:], decodedPK[:])

	//Complete the post-signup steps, including challenge generation and issuance of a temporary token
	if err := postSignup(w, r, user, userCollection); err != nil {
		httpu.HttpErrorAsJson(w, err, http.StatusInternalServerError)
		return
	}
}

/*
Ensures that a user doesn't already exist in the database based on what
was given by the user. A `nil` error indicates that no matching records
were found. Checking collections for existent objects is expensive, so
not all records are checked if one fails.
*/
func ensureNonexistantUser(coll *mongo.Collection, usr intermediateUser, ctx context.Context) error {
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
	hits, err := mongoutil.Aggregate(coll, agg, ctx)
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
func postSignup(w http.ResponseWriter, r *http.Request, user *obj.User, ucoll *mongo.Collection) error {
	//Step 1: Issue a token that's good for the duration of the challenge window; otherwise the routes won't be allowed
	tempToken := obj.NewToken(user.ID, ip_addr.HttpIP2NetIP(r.RemoteAddr), obj.TokenScopePOSTSIGNUP, user.Flags.PurgeBy)
	fmt.Printf("TOK: '%s'\n", tempToken.ToB64())

	//Step 2a: Push the token to the user's list of tokens and add the user to the database
	//TODO: use CRUD operations here
	user.Tokens = append(user.Tokens, *tempToken)
	userBson, jerr := bson.Marshal(user)
	_, ierr := ucoll.InsertOne(r.Context(), userBson)
	if jerr != nil {
		return jerr
	}
	if ierr != nil {
		return ierr
	}

	//Set 2b: Cache the access tokens
	cr.CreateSA(rcl, r.Context(), user.ID.UUID, tempToken.String())

	//Step 3a: Create challenges for email and public key verification
	srvIdent := obj.Identifiable{ID: env.ID, Type: obj.IdTypeSERVER}
	usrIdent := obj.Identifiable{ID: user.ID, Type: user.Type}
	expiry := user.Flags.PurgeBy
	emailChall := challenge.NewChallenge(challenge.ChallengeScopeEMAIL, srvIdent, usrIdent, expiry)
	pubkeyChall := challenge.NewChallenge(challenge.ChallengeScopePUBKEY, srvIdent, usrIdent, expiry)

	//Step 3b: Compose the challenge URL for the email
	baseUrl := "http://127.0.0.1:8888" //TODO: change this eventually
	echallUrl := util.Must(url.Parse(fmt.Sprintf("%s/challenges/%s/solve", baseUrl, emailChall.ID)))
	eurlParams := echallUrl.Query()
	eurlParams.Set(mw.AuthHttpParamName, tempToken.ToB64())
	eurlParams.Set(challenge.ChallengeURLParamName, emailChall.Payload)
	echallUrl.RawQuery = eurlParams.Encode()

	//Step 3c: Compose the challenge email to send to the user
	emsg := mail.NewMSG()
	emsg.SetFrom(cfg.Email.Username)
	emsg.AddTo(user.Email)
	emsg.SetSubject("Your Wraith Account")

	//Step 3d: Create the body of the email
	tmplFields := remailt.Template{
		UUID:          user.ID.String(),
		UName:         user.Username,
		Email:         user.Email,
		PKFingerprint: user.Pubkey.Fingerprint(),
		PurgeTime:     util.Time2OffsetReq(user.Flags.PurgeBy, r).Format(time.RFC1123Z),
		ChallengeLink: echallUrl.String(),
	}
	var ebody bytes.Buffer
	if err := emailChallTemplate.Execute(&ebody, tmplFields); err != nil {
		return err
	}
	emsg.SetBody(mail.TextHTML, ebody.String())

	//Step 3e: Send the email challenge to the user's email
	if emsg.Error != nil {
		return emsg.Error
	}
	if err := email.GetInstance().SendEmail(emsg); err != nil {
		return err
	}

	//Step 4: Push the challenges to the database for later retrieval
	crud.AddChallenges(mcl, rcl, r.Context(), emailChall, pubkeyChall)

	//Step 5: Write the response back to the user
	psu := postsignupUser{
		ID:              user.ID,
		RedactedEmail:   util.RedactEmail(user.Email),
		Challenges:      []mongoutil.UUID{emailChall.ID, pubkeyChall.ID},
		TempAccessToken: tempToken.ToB64(),
	}
	if jerr := json.NewEncoder(w).Encode(&psu); jerr != nil {
		return jerr
	}

	//No errors so return nil
	return nil
}
