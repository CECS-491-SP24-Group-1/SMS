package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tanqiangyes/govalidator"
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

// Validates an `intermediateUser` object using `tanqiangyes/govalidator`.
func (iu intermediateUser) validate(strictEmail bool) (bool, []error) {
	//Create a slice to hold the collected errors
	errors := []error{}

	//Step 1: Check the validity of the username
	//Username should be 4-16 characters in length and only consist of alphanumeric characters and underscores
	//This function should never throw an error
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

	//Do something with the object
	userStr := fmt.Sprintf("User: %+v", user)
	fmt.Printf("%s\n", userStr)

	//Respond back with the UUID of the new user
	if err := json.NewEncoder(w).Encode(user); err != nil {
		util.HttpErrorAsJson(w, err, http.StatusInternalServerError)
	}
}
