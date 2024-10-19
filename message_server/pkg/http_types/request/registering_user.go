package request

import (
	"fmt"
	"regexp"
	"strings"

	"wraith.me/message_server/pkg/controller/csolver"
	"wraith.me/message_server/pkg/crypto"
)

/*
Represents a user object sent during a user registration request.
*/
type RegisteringUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pubkey   string `json:"pubkey"`
}

// Validates a `RegisteringUser` object using `tanqiangyes/govalidator`.
func (ru RegisteringUser) Validate(strictEmail bool) (bool, []error) {
	//Create a slice to hold the collected errors
	errors := []error{}

	//Step 1: Check the validity of the username
	//Username should be 4-16 characters in length and only consist of alphanumeric characters and underscores
	//This function should never throw an error since the regexp is hard-coded
	validUname, _ := regexp.MatchString(`^([a-z0-9_]){4,16}$`, strings.ToLower(ru.Username))
	if !validUname {
		errors = append(errors, fmt.Errorf("username '%s' is invalid. It must be 4-16 characters in length and only consist of alphanumeric characters and underscores", strings.ToLower(ru.Username)))
	}

	//Step 2: Check the validity of the email
	evErr := csolver.IsValidEmail(ru.Email, strictEmail)
	if evErr != nil {
		errors = append(errors, evErr)
	}

	//Step 3: Check the validity of the base64'ed public key by attempting to convert to a byte array
	validPubkey := true
	_, err := crypto.ParsePubkey(ru.Pubkey)
	if err != nil {
		validPubkey = false
		errors = append(errors, err)
	}

	//Return the validity status and any errors that occurred
	return validUname && (evErr == nil) && validPubkey, errors
}
