package token

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	ccrypto "wraith.me/message_server/crypto"
	"wraith.me/message_server/util"
)

const (
	_TOK_TYPE = "ttype"
	_TOK_IP   = "tipaddr"
)

var (
	//The name of the access token cookie.
	AccessTokenName = "access_token"

	//The name of the access token expiration cookie.
	AccessTokenExprName = "access_token_expr"

	//Controls whether to add the token expiry to the footer.
	ExprInFooter = true

	//The name of the refresh token cookie.
	RefreshTokenName = "refresh_token"

	//The name of the refresh token expiration cookie.
	RefreshTokenExprName = "refresh_token_expr"

	//The format of the times in the token.
	TimeFmt = time.RFC3339
)

//
//-- CLASS: Token
//

/*
Represents a PASETO token that allows users to authenticate with the API
in a semi-stateless manner after login. This class can represent either a
fully stateless, short-lived access token or a stateful, long-refresh token.
*/
type Token struct {
	//The ID of the token. This is the `jti` field of the PASETO token. This is calculated from the `iat` field.
	ID util.UUID `json:"id"`

	//The ID of the entity that issued the token. This is the `iss` field of the PASETO token.
	Issuer util.UUID `json:"issuer"`

	//The user that this token is for by ID. This is the `sub` field of the PASETO token.
	Subject util.UUID `json:"subject"`

	//The time at which the token should expire. This is the `exp` field of the PASETO token.
	Expiry time.Time `json:"expires"`

	//The time at which the token was issued. This is the `iat` and `nbf` fields of the PASETO token.
	Issued time.Time `json:"issued"`

	//The type of token this is.
	Type TokenType `json:"type"`

	//The IP address of the client that the token was originally created for
	IPAddr net.IP `json:"ip_addr"`
}

//-- Constructors

/*
Constructs a new token object, which takes in the subject, issuer, type,
expiry, and optional time to use for the `iat` and `nbf` fields.
*/
func NewToken(subject util.UUID, issuer util.UUID, typ TokenType, exp time.Time, now *time.Time) *Token {
	//Check if the "now" parameter is nil
	if now == nil {
		n := time.Now()
		now = &n
	}

	//Generate a `jti` value based on the "now" time
	jti := util.NewUUID7FromTime(*now)

	//Construct the object
	//IP is added later
	return &Token{
		ID:      jti,
		Issuer:  issuer,
		Subject: subject,
		Expiry:  exp,
		Issued:  *now,
		Type:    typ,
	}
}

//-- Methods

// Creates an HTTP cookie string to hold the PASETO token.
func (t Token) Cookie(key ccrypto.Privkey, path, domain string, persistent bool) http.Cookie {
	_, cookie := t.CryptAndCookie(key, path, domain, persistent)
	return cookie
}

/*
Encrypts the PASETO token and creates an HTTP cookie string to hold the
PASETO token, all in one step.
*/
func (t Token) CryptAndCookie(key ccrypto.Privkey, path, domain string, persistent bool) (token string, cookie http.Cookie) {
	//Encrypt the token
	token = t.Encrypt(key, ExprInFooter)

	//Get the name based on whether this is an access or refresh token
	name := "Untitled"
	if t.Type == TokenTypeACCESS {
		name = AccessTokenName
	} else if t.Type == TokenTypeREFRESH {
		name = RefreshTokenName
	}

	//Build the cookie
	cookieBuilder := http.Cookie{
		Name:     name,
		Value:    token,
		Path:     path,
		Domain:   domain,
		MaxAge:   t.MaxAge(persistent),
		Secure:   true,
		HttpOnly: true, //This must be true to ensure it remains inaccessible by clientside JS
		SameSite: http.SameSiteStrictMode,
	}
	cookie = cookieBuilder
	return
}

// Decrypts a PASETO string using a given symmetric key, which creates a Token object.
func Decrypt(token string, key ccrypto.Privkey, issuer util.UUID, typ TokenType) (*Token, error) {
	//Create a new token parser and add basic rules
	parser := paseto.NewParser()
	parser.AddRule(
		paseto.ValidAt(time.Now()),       //Checks nbf, iat, and exp in one fell-swoop
		paseto.IssuedBy(issuer.String()), //Ensures this server issued the token
		matchingType(typ),                //Token type and input purpose must match
	)

	//Decrypt the token and validate it; due to the "v4_local" construction, any tamper attempts will auto-fail this check
	decrypted, err := parser.ParseV4Local(util.Edsk2PasetoSK(key), token, nil)
	if err != nil {
		return nil, err
	}

	//Decode the token and return the payload
	return pasetoDecode(decrypted, issuer)
}

/*
Encrypts this token using a given symmetric key, which creates a PASETO
token string. Optionally, the token can include the expiration in the
footer.
*/
func (t Token) Encrypt(key ccrypto.Privkey, expInFooter bool) string {
	//Create a new token with expiration in x time
	token := paseto.NewToken()
	token.SetIssuedAt(t.Issued)   //Token "iat"
	token.SetNotBefore(t.Issued)  //Token "nbf"
	token.SetExpiration(t.Expiry) //Token "exp"

	//Add additional data to the token
	token.SetJti(t.ID.String())                 //Token ID
	token.SetIssuer(t.Issuer.String())          //Issuer ID (server)
	token.SetSubject(t.Subject.String())        //User ID (client)
	token.SetString(_TOK_TYPE, t.Type.String()) //Token type
	token.SetString(_TOK_IP, t.IPAddr.String()) //Subject IP

	//Check if the expiration footer should be added
	if expInFooter {
		token.SetFooter([]byte(t.Expiry.UTC().Format(TimeFmt)))
	}

	//Encrypt the token
	return token.V4Encrypt(util.Edsk2PasetoSK(key), nil)
}

// Creates an HTTP cookie string to hold the expiration of this token.
func (t Token) ExprCookie(path, domain string, exprMultiplier int, persistent bool) http.Cookie {
	//Get the name based on whether this is an access or refresh token
	name := "Untitled"
	if t.Type == TokenTypeACCESS {
		name = AccessTokenExprName
	} else if t.Type == TokenTypeREFRESH {
		name = RefreshTokenExprName
	}

	//Build the cookie
	cookieBuilder := http.Cookie{
		Name:     name,
		Value:    t.Expiry.Format(TimeFmt),
		Path:     path,
		Domain:   domain,
		MaxAge:   t.MaxAge(persistent) * exprMultiplier,
		Secure:   true,
		HttpOnly: false, //This must be false so clientside JS can access it
		SameSite: http.SameSiteStrictMode,
	}
	return cookieBuilder
}

// Gets the "Max Age" of the token in seconds.
func (t Token) MaxAge(persistent bool) int {
	// Get the delta between the current date and the expiry; in seconds
	timeDelta := -1
	if persistent {
		timeDelta = int(math.Round(time.Until(t.Expiry).Seconds()))
	}
	return timeDelta
}

//-- Public utilities

// Gets the expiration from a token that has it in the footer
func GetExprFromFooter(tok string) (time.Time, error) {
	//Split the token at every period and get the last piece
	pieces := strings.Split(tok, ".")
	exprB64 := pieces[len(pieces)-1]

	/*
		Ensure this token has a timestamp; expected size of the split pieces is 4:
		- v4
		- local
		- <token>
		- <footer>
	*/
	const expectedSize = 4
	if len(pieces) != expectedSize {
		return time.Time{}, fmt.Errorf("token doesn't have a valid footer; got size %d expected %d", len(pieces), expectedSize)
	}

	//Decode the expiration from base64
	exprBytes, err := base64.RawURLEncoding.DecodeString(exprB64)
	if err != nil {
		return time.Time{}, err
	}

	//Parse the timestamp to a `time.Time`
	return time.Parse(TimeFmt, string(exprBytes))
}

//-- Private utilities

// Decodes a PasetoV4 token into a valid `Token` object.
func pasetoDecode(tok *paseto.Token, issuer util.UUID) (*Token, error) {
	//Get the fields of the token
	var id string
	//var issuer string
	var subject string
	var expiry time.Time
	var issued time.Time
	var typ string
	var ipAddr string

	//Early return if any conversion function fails
	//TODO: might want to condense this
	perr := func() (err error) {
		id, err = tok.GetJti()
		if err != nil {
			return
		}
		subject, err = tok.GetSubject()
		if err != nil {
			return
		}
		expiry, err = tok.GetExpiration()
		if err != nil {
			return
		}
		issued, err = tok.GetIssuedAt()
		if err != nil {
			return
		}
		typ, err = tok.GetString(_TOK_TYPE)
		if err != nil {
			return
		}

		ipAddr, err = tok.GetString(_TOK_IP)
		if err != nil {
			return
		}
		return nil
	}()
	if perr != nil {
		return nil, perr
	}

	//Create a new struct and return it
	return &Token{
		ID:      util.UUIDFromString(id),
		Issuer:  issuer,
		Subject: util.UUIDFromString(subject),
		Expiry:  expiry,
		Issued:  issued,
		Type:    MustParseTokenType(typ),
		IPAddr:  net.ParseIP(ipAddr),
	}, nil
}

// Verifies that the token's purpose matches an input one.
func matchingType(typ TokenType) paseto.Rule {
	return func(token paseto.Token) error {
		//Get the token type from the token
		ttype, err := token.GetString(_TOK_TYPE)
		if err != nil {
			return err
		}

		//Parse the token type to a string
		typeo, err := ParseTokenType(ttype)
		if err != nil {
			return err
		}

		//Check if the token purpose is appropriate
		if subtle.ConstantTimeByteEq(uint8(typeo), uint8(typ)) == 0 {
			return fmt.Errorf("this token's type is not appropriate; must be '%s'", typ.String())
		}

		//No error, so return `nil`
		return nil
	}
}
