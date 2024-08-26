package csolver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"wraith.me/message_server/config"
	"wraith.me/message_server/obj/challenge"
	c "wraith.me/message_server/obj/challenge"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/template/reg_email"
	"wraith.me/message_server/util"
)

// Issues an email challenge for a user. This is stage 1 of an email challenge.
func IssueEmailChallenge(usr *user.User, cfg *config.Config, env *config.Env, r *http.Request) error {
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

	//Nothing went wrong, so return nil
	return nil
}

// Verifies that an email challenge is valid. This is stage 2 of an email challenge.
func VerifyEmailChallenge(env *config.Env, ctext string, w http.ResponseWriter, r *http.Request) *c.CToken {
	//Get the challenge text
	if strings.TrimSpace(ctext) == "" {
		//Bail out if nothing was supplied
		util.ErrResponse(
			http.StatusBadRequest,
			fmt.Errorf("received empty challenge response"),
		).Respond(w)
		return nil
	}

	//Attempt to decrypt the challenge
	//From this point on, it's safe to assume the user successfully passed the challenge
	ctoken, err := challenge.Decrypt(ctext, env.SK, env.ID, challenge.CPurposeCONFIRM)
	if err != nil {
		util.ErrResponse(http.StatusForbidden, err).Respond(w)
		return nil
	}

	//Return the token
	return ctoken
}

// Checks if a user's email is valid
func IsValidEmail(email string, strict bool) error {
	//`strictEmail` also ensures the email maps to an existing domain name
	emailValidator := govalidator.IsEmail
	if strict {
		emailValidator = govalidator.IsExistingEmail
	}

	//Step 2: Check the validity of the email
	validEmail := emailValidator(strings.ToLower(email))
	if !validEmail {
		return fmt.Errorf("email '%s' is invalid; it must be of the form 'foo@example.com'", strings.ToLower(email))
	}

	//Nothing went wrong, so return nil
	return nil
}
