package tests

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"wraith.me/message_server/pkg/util"
)

const SUBJECT_PK_KEY = "sub-pk"

// See: https://developer.okta.com/blog/2019/10/17/a-thorough-introduction-to-paseto
func TestPaseto(t *testing.T) {
	//Create a new token with expiration in 10 minutes
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(10 * time.Minute))

	//Set additional fields for the token
	jtiIn := util.MustNewUUID7().String()
	issIn := util.MustNewUUID7().String()
	subIn := util.MustNewUUID7().String()
	subpkIn := base64.StdEncoding.EncodeToString(paseto.NewV4SymmetricKey().ExportBytes())

	//Add additional data to the token
	token.SetJti(jtiIn)                      //Token ID
	token.SetIssuer(issIn)                   //Issuer ID (server)
	token.SetSubject(subIn)                  //User ID (client)
	token.SetString(SUBJECT_PK_KEY, subpkIn) //User public key

	//Create a new symmetric key and encrypt the token
	key := paseto.NewV4SymmetricKey() // Don't share this!! This can be generated once and used for as many tokens as required
	encrypted := token.V4Encrypt(key, nil)

	//Print the token to the console
	fmt.Println("-------------- I N --------------")
	fmt.Printf("TKS: `%s`\n", encrypted)
	printTokInfo(key, token)

	//
	//------------------------------------------------------------------------------
	//

	//Create a new parser and add rules
	parser := paseto.NewParser()
	parser.AddRule(paseto.ValidAt(time.Now())) //Checks nbf, iat, and exp in one fell-swoop
	parser.AddRule(paseto.IssuedBy(issIn))
	parser.AddRule(paseto.Subject(subIn))
	parser.AddRule(subjectHasPK(subpkIn))

	//Decrypt the token and validate it; due to the "v4_local" construction, any tamper attempts will auto-fail this check
	decrypted, err := parser.ParseV4Local(key, encrypted, nil)
	if err != nil {
		t.Fatal(err)
	}

	//Print the token info
	fmt.Println("-------------- OUT --------------")
	printTokInfo(key, *decrypted)

	//
	//------------------------------------------------------------------------------
	//

	//Purposefully tamper with the token at a random byte
	erunes := []rune(encrypted)
	erunes[rand.Intn(len(erunes))] = '\x69'
	tampered := string(erunes)

	//Ensure the tampered token is not accepted
	_, err = parser.ParseV4Local(key, tampered, nil)
	if err != nil {
		t.Logf("Tampered token was successfully rejected with message `%s`", err)
	} else {
		t.Fatalf("tampered token was not rejected!")
	}
}

// Verifies that the public key of the subject matches an input value.
func subjectHasPK(pk string) paseto.Rule {
	return func(token paseto.Token) error {
		//Get the public key from the token
		tpk, err := token.GetString(SUBJECT_PK_KEY)
		if err != nil {
			return err
		}

		//Check the validity of the subject's public key using constant time compare
		//This prevents side channel attacks against this field of the token
		if subtle.ConstantTimeCompare([]byte(tpk), []byte(pk)) == 0 {
			return fmt.Errorf("this token's subject has no public key `%s'. `%s' found", pk, tpk)
		}

		//No error, so return `nil`
		return nil
	}
}

func printTokInfo(key paseto.V4SymmetricKey, tok paseto.Token) {
	//fmt.Println("-------------- TOK --------------")
	fmt.Printf("Key:           %s\n", base64.StdEncoding.EncodeToString(key.ExportBytes()))
	jti, _ := tok.GetJti()
	iss, _ := tok.GetIssuer()
	sub, _ := tok.GetSubject()
	subpk, _ := tok.GetString(SUBJECT_PK_KEY)
	fmt.Printf("Token ID:      %s\n", jti)
	fmt.Printf("Token Issuer:  %s\n", iss)
	fmt.Printf("Token Subject: %s\n", sub)
	fmt.Printf("Token Sub PK:  %s\n", subpk)
	fmt.Printf("Token Full:    %s\n", tok)
	fmt.Println("-------------- +++ --------------")
}
