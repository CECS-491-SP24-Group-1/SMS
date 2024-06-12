package challenge

import (
	"crypto/subtle"
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	ccrypto "wraith.me/message_server/crypto"
	"wraith.me/message_server/db/mongoutil"
)

const (
	_CHALL_CTYPE    = "ctype"
	_CHALL_CPURPOSE = "cpurpose"
	_CHALL_CLAIM    = "claim"
)

/*
Represents a challenge given to a user to solve. A challenge can be used
to remove holds on accounts, prove identity, or provide authorization for
an account action such as deletion. A challenge can either be initiated by
a user or a server. Likewise, a challenge can either be responded to by a
user or a server, though the latter is not currently slated for immediate
implementation at this time. This implementation of a challenge is meant
to be encoded in a stateless PASETO token. This token is then echoed back
to the server in the case of an email challenge and echoed plus signed in
the case of a public key challenge. The `Claim` field can either represent
an email or a base64-encoded public key, and the state depends on the value
of the `CType`field.
*/
type CToken struct {
	ID        mongoutil.UUID
	Issuer    mongoutil.UUID
	SubjectID mongoutil.UUID
	CType     CType
	Purpose   CPurpose
	Expiry    time.Time
	Claim     string
}

// Creates a new challenge meant for validating ownership of an email.
func NewEmailChallenge(issuer mongoutil.UUID, subjectID mongoutil.UUID, purpose CPurpose, expiry time.Time, email string) CToken {
	return CToken{
		ID:        mongoutil.MustNewUUID7(),
		Issuer:    issuer,
		SubjectID: subjectID,
		CType:     CTypeEMAIL,
		Purpose:   purpose,
		Expiry:    expiry,
		Claim:     email,
	}
}

// Creates a new challenge meant for validating ownership of a private key.
func NewPKChallenge(issuer mongoutil.UUID, subjectID mongoutil.UUID, purpose CPurpose, expiry time.Time, pubkey ccrypto.Pubkey) CToken {
	return CToken{
		ID:        mongoutil.MustNewUUID7(),
		Issuer:    issuer,
		SubjectID: subjectID,
		CType:     CTypePUBKEY,
		Purpose:   purpose,
		Expiry:    expiry,
		Claim:     pubkey.String(),
	}
}

// Encodes a challenge payload into an encrypted v4 Paseto token.
func (t CToken) Encrypt(key ccrypto.Privkey) string {
	//Create a new token with expiration in 10 minutes
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())                         //Token "iat"
	token.SetNotBefore(time.Now())                        //Token "nbf"
	token.SetExpiration(time.Now().Add(10 * time.Minute)) //Token "exp"

	//Add additional data to the token
	token.SetJti(t.ID.String())                          //Token ID
	token.SetIssuer(t.Issuer.String())                   //Issuer ID (server)
	token.SetSubject(t.SubjectID.String())               //User ID (client)
	token.SetString(_CHALL_CTYPE, t.CType.String())      //Challenge type
	token.SetString(_CHALL_CPURPOSE, t.Purpose.String()) //Challenge purpose
	token.SetString(_CHALL_CLAIM, t.Claim)               //Challenge claim

	//Encrypt the token
	return token.V4Encrypt(edsk2PasetoSK(key), nil)
}

// Decodes an encrypted v4 Paseto token into a challenge payload.
func Decrypt(token string, key ccrypto.Privkey, issuer mongoutil.UUID, purpose CPurpose) (*CToken, error) {
	return decryptBackend(token, key, issuer, purpose)
}

// Decodes an encrypted v4 Paseto token into a challenge payload, with stricter checks.
func DecryptPKStrict(token string, key ccrypto.Privkey, issuer mongoutil.UUID, purpose CPurpose, subject mongoutil.UUID, pubkey ccrypto.Pubkey) (*CToken, error) {
	//Create a list of additional rules
	rules := []paseto.Rule{
		paseto.Subject(subject.String()), //Token subject and input subject must match
		subjectHasPK(pubkey.String()),    //Token subject PK and input subject PK must match
	}

	//Decrypt the token and add the extra rules
	return decryptBackend(token, key, issuer, purpose, rules...)
}

// Performs token decryption and ensures it matches against a rule-set.
func decryptBackend(token string, key ccrypto.Privkey, issuer mongoutil.UUID, purpose CPurpose, rules ...paseto.Rule) (*CToken, error) {
	//Create a new token parser and add basic rules
	parser := paseto.NewParser()
	parser.AddRule(
		paseto.ValidAt(time.Now()),       //Checks nbf, iat, and exp in one fell-swoop
		paseto.IssuedBy(issuer.String()), //Ensures this server issued the token
		matchingPurpose(purpose),         //Token purpose and input purpose must match
	)

	//Add additional rules from the rules array
	parser.AddRule(rules...)

	//Decrypt the token and validate it; due to the "v4_local" construction, any tamper attempts will auto-fail this check
	decrypted, err := parser.ParseV4Local(edsk2PasetoSK(key), token, nil)
	if err != nil {
		return nil, err
	}

	//Decode the token and return the payload
	return pasetoDecode(decrypted, issuer)
}

// Converts an Ed25519 SK to a PasetoV4 SK.
func edsk2PasetoSK(key ccrypto.Privkey) paseto.V4SymmetricKey {
	seed := key.Seed()
	psk, _ := paseto.V4SymmetricKeyFromBytes(seed[:])
	return psk
}

// Decodes a PasetoV4 token into a valid `CToken` object.
func pasetoDecode(tok *paseto.Token, issuer mongoutil.UUID) (*CToken, error) {
	//Get the fields of the token
	var id string
	var subjectID string
	var ctype CType
	var purpose CPurpose
	var expiry time.Time
	var claim string

	//Early return if any conversion function fails
	perr := func() (err error) {
		id, err = tok.GetJti()
		if err != nil {
			return
		}
		subjectID, err = tok.GetSubject()
		if err != nil {
			return
		}
		ctypeS, err := tok.GetString(_CHALL_CTYPE)
		if err != nil {
			return
		}
		ctype, err = ParseCType(ctypeS)
		if err != nil {
			return
		}
		purposeS, err := tok.GetString(_CHALL_CPURPOSE)
		if err != nil {
			return
		}
		purpose, err = ParseCPurpose(purposeS)
		if err != nil {
			return
		}
		expiry, err = tok.GetExpiration()
		if err != nil {
			return
		}
		claim, err = tok.GetString(_CHALL_CLAIM)
		if err != nil {
			return
		}
		return nil
	}()
	if perr != nil {
		return nil, perr
	}

	//Ensure the claim maps to a valid public key if it is one
	if ctype == CTypePUBKEY {
		if _, err := ccrypto.ParsePubkey(claim); err != nil {
			return nil, err
		}
	}

	//Create a new struct and return it
	return &CToken{
		ID:        mongoutil.UUIDFromString(id),
		Issuer:    issuer,
		SubjectID: mongoutil.UUIDFromString(subjectID),
		CType:     ctype,
		Purpose:   purpose,
		Expiry:    expiry,
		Claim:     claim,
	}, nil
}

// Verifies that the public key of the subject matches an input value.
func subjectHasPK(pk string) paseto.Rule {
	return func(token paseto.Token) error {
		//Get the public key from the token
		tpk, err := token.GetString(_CHALL_CLAIM)
		if err != nil {
			return err
		}

		//Get the challenge type from the token
		ttypes, err := token.GetString(_CHALL_CTYPE)
		if err != nil {
			return err
		}

		//Parse the challenge type to a string
		ttype, err := ParseCType(ttypes)
		if err != nil {
			return err
		}

		//Check if the token type is appropriate
		if subtle.ConstantTimeByteEq(uint8(ttype), uint8(CTypePUBKEY)) == 0 {
			return fmt.Errorf("this token's type is not appropriate; must be 'pubkey'")
		}

		//Check the validity of the subject's public key using constant time compare
		//This prevents side channel attacks against this field of the token
		if subtle.ConstantTimeCompare([]byte(tpk), []byte(pk)) == 0 {
			return fmt.Errorf("this token's subject has a mismatched or no public key")
		}

		//No error, so return `nil`
		return nil
	}
}

// Verifies that the token's purpose matches an input one.
func matchingPurpose(purpose CPurpose) paseto.Rule {
	return func(token paseto.Token) error {
		//Get the challenge purpose from the token
		tpups, err := token.GetString(_CHALL_CPURPOSE)
		if err != nil {
			return err
		}

		//Parse the challenge purpose to a string
		tpup, err := ParseCPurpose(tpups)
		if err != nil {
			return err
		}

		//Check if the token purpose is appropriate
		if subtle.ConstantTimeByteEq(uint8(tpup), uint8(purpose)) == 0 {
			return fmt.Errorf("this token's purpose is not appropriate; must be '%s'", purpose.String())
		}

		//No error, so return `nil`
		return nil
	}
}
