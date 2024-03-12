package users

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tanqiangyes/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/db"
	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
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
Ensures that a user doesn't already exist in the database based on what
was given by the user. A `nil` error indicates that no matching records
were found. Checking collections for existant objects is expensive, so
not all records are checked if one fails.
*/
func ensureNonexistantUser(coll *mongo.Collection, usr intermediateUser, ctx context.Context) error {
	//Parse out the public key of the incoming user
	pubkey, _ := obj.ParsePubkeyBytes(usr.Pubkey) //Errors should not occur here; data is already pre-validated

	//Construct a Mongo aggregation pipeline to run the request; avoids making multiple round-trips to the database
	//This aggregation was exported from MongoDB; do not edit if you don't know what you are doing!
	agg := bson.A{
		bson.D{
			//Match any documents that have the same username, email, or public key
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
	dbytes, err := base64.StdEncoding.DecodeString(iu.Pubkey)
	if err != nil {
		errors = append(errors, err)
	}
	validPubkey := len(dbytes) == obj.PUBKEY_SIZE
	if !validPubkey {
		errors = append(errors, fmt.Errorf("mismatched public key size (%d); expected: %d", len(dbytes), obj.PUBKEY_SIZE))
	}

	//Return the validity status and any errors that occurred
	return validUname && validEmail && validPubkey, errors
}

// Handles incoming requests made to `POST /users/register`.
func RegisterUserRoute(w http.ResponseWriter, r *http.Request) {
	//Create a new intermediate user object
	iuser := intermediateUser{}

	//Get the request body and attempt to parse to JSON
	if err := json.NewDecoder(r.Body).Decode(&iuser); err != nil {
		util.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Validate the unmarshalled struct
	//At this point, the incoming JSON was accepted, but fields may be missing or invalid
	valid, verrs := iuser.validate(false)
	if !valid {
		util.HttpMultipleErrorsAsJson(w, verrs, http.StatusBadRequest)
		return
	}

	//Decode the base64 public key to a byte array
	decodedPK, err := base64.StdEncoding.DecodeString(iuser.Pubkey)
	if err != nil {
		util.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Get the users collection from the database and ensure a record doesn't already exist
	dbc := db.GetInstance().GetClient()
	userCollection := dbc.Database(db.ROOT_DB).Collection(db.USERS_COLLECTION)

	//Ensure the user doesn't already exist in the database
	if err := ensureNonexistantUser(userCollection, iuser, r.Context()); err != nil {
		util.HttpErrorAsJson(w, err, http.StatusBadRequest)
		return
	}

	//Fill in the rest of the details
	uuid, _ := mongoutil.NewUUID7()
	user := obj.User{
		ID:          *uuid,
		Username:    strings.ToLower(iuser.Username),
		DisplayName: iuser.Username,
		Email:       strings.ToLower(iuser.Email),
		LastLogin:   time.Now(),
		LastIP:      util.HttpIP2NetIP(r.RemoteAddr),
		Flags:       obj.DefaultUserFlags(),
	}
	copy(user.Pubkey[:], decodedPK[:])

	//Add the object to the database
	userBson, jerr := bson.Marshal(user)
	_, ierr := userCollection.InsertOne(r.Context(), userBson)
	if jerr != nil {
		fmt.Printf("JERR: %s\n", jerr.Error())
		util.HttpErrorAsJson(w, jerr, http.StatusInternalServerError)
		return
	}
	if ierr != nil {
		fmt.Printf("IERR: %s\n", ierr.Error())
		util.HttpErrorAsJson(w, ierr, http.StatusInternalServerError)
		return
	}

	//Do something with the object
	userStr := fmt.Sprintf("User: %+v", user)
	fmt.Printf("%s\n", userStr)

	//Respond back with the UUID of the new user
	if err := json.NewEncoder(w).Encode(user); err != nil {
		util.HttpErrorAsJson(w, err, http.StatusInternalServerError)
	}
}
