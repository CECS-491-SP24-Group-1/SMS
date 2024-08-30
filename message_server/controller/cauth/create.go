package cauth

import (
	"context"
	"net/http"
	"time"

	"wraith.me/message_server/config"
	"wraith.me/message_server/obj/ip_addr"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util"
)

var (
	//Controls the path used for expiration cookies.
	ExprCookiePath = "/"
)

/*
Issues an access token for a user and writes it to the outgoing response
cookies.
*/
func IssueAccessToken(w http.ResponseWriter, r *http.Request, usr *user.User, env *config.Env, cfg *token.TConfig, parent *util.UUID, persistent bool) {
	//Get the current and expiry times
	now := time.Now()
	exp := now.Add(time.Duration(cfg.AccessLifetime) * time.Second)

	//Create an access token object and add additional fields
	atoken := token.NewToken(
		usr.ID,
		env.ID,
		token.TokenTypeACCESS,
		exp,
		parent,
		&now,
	)
	atoken.IPAddr = ip_addr.HttpIP2NetIP(r.RemoteAddr)
	atoken.UserAgent = r.UserAgent()

	//Encrypt the access token and generate the cookie strings
	domain := cfg.Domain
	atCookie := atoken.Cookie(
		env.SK, cfg.AccessCookiePath,
		domain, persistent,
	)
	ateCookie := atoken.ExprCookie(
		ExprCookiePath, domain,
		cfg.ExprMultiplier, persistent,
	)

	//Write the cookies to the outgoing response
	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &ateCookie)
}

/*
Issues a refresh token for a user and writes it to the outgoing response
cookies along with the user object in the database. It is assumed that the
user in question already exists in the database.
*/
func IssueRefreshToken(w http.ResponseWriter, r *http.Request, usr *user.User, ucoll *user.UserCollection, ctx context.Context, env *config.Env, cfg *token.TConfig, persistent bool) (util.UUID, error) {
	//Get the current and expiry times
	now := time.Now()
	exp := now.Add(time.Duration(cfg.RefreshLifetime) * time.Second)

	//Create a refresh token object and add additional fields
	rtoken := token.NewToken(
		usr.ID,
		env.ID,
		token.TokenTypeREFRESH,
		exp,
		nil,
		&now,
	)
	rtoken.IPAddr = ip_addr.HttpIP2NetIP(r.RemoteAddr)
	rtoken.UserAgent = r.UserAgent()

	//Encrypt the refresh token and generate the cookie strings
	domain := cfg.Domain
	rte, rtCookie := rtoken.CryptAndCookie(
		env.SK, cfg.RefreshCookiePath,
		domain, persistent,
	)
	rteCookie := rtoken.ExprCookie(
		ExprCookiePath, domain,
		cfg.ExprMultiplier, persistent,
	)

	//Write the cookies to the outgoing response
	http.SetCookie(w, &rtCookie)
	http.SetCookie(w, &rteCookie)

	//Add the refresh token to the user's list of tokens, keyed by its ID
	usr.AddToken(rtoken.ID.String(), rte, rtoken.Expiry)

	//Upsert the corresponding document in the database
	_, err := ucoll.UpsertId(ctx, usr.ID, usr)
	return rtoken.ID, err
}
