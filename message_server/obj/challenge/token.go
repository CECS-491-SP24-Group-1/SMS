package challenge

import (
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
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(10 * time.Minute))

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
func Decrypt(token string, issuer mongoutil.UUID, key ccrypto.Privkey) (*CToken, error) {
	//Create a new token parser and add rules
	parser := paseto.NewParser()
	parser.AddRule(paseto.ValidAt(time.Now())) //Checks nbf, iat, and exp in one fell-swoop
	parser.AddRule(paseto.IssuedBy(issuer.String()))
	//parser.AddRule(paseto.Subject(subIn))

	//Decrypt the token and validate it; due to the "v4_local" construction, any tamper attempts will auto-fail this check
	decrypted, err := parser.ParseV4Local(edsk2PasetoSK(key), token, nil)
	if err != nil {
		return nil, err
	}

	//Get the fields of the token
	var id string
	var subjectID string
	var ctype CType
	var purpose CPurpose
	var expiry time.Time
	var claim string

	//Early return if any conversion function fails
	perr := func() (err error) {
		id, err = decrypted.GetJti()
		if err != nil {
			return
		}
		subjectID, err = decrypted.GetSubject()
		if err != nil {
			return
		}
		ctypeS, err := decrypted.GetString(_CHALL_CTYPE)
		if err != nil {
			return
		}
		ctype, err = ParseCType(ctypeS)
		if err != nil {
			return
		}
		purposeS, err := decrypted.GetString(_CHALL_CPURPOSE)
		if err != nil {
			return
		}
		purpose, err = ParseCPurpose(purposeS)
		if err != nil {
			return
		}
		expiry, err = decrypted.GetExpiration()
		if err != nil {
			return
		}
		claim, err = decrypted.GetString(_CHALL_CLAIM)
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

// Converts an Ed25519 SK to a PasetoV4 SK.
func edsk2PasetoSK(key ccrypto.Privkey) paseto.V4SymmetricKey {
	seed := key.Seed()
	psk, _ := paseto.V4SymmetricKeyFromBytes(seed[:])
	return psk
}
