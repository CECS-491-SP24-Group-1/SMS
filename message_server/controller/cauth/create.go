package cauth

const (
	//Controls whether to add the token expiry to the footer.
	AddExpToFooter = true
)

/*
Issues an access token for a user and writes it to the outgoing response
cookies.
* /
func IssueAccessToken(w http.ResponseWriter, usr *user.User, env *config.Env, cfg *token.TConfig, persistent bool) {
	//Get the current and expiry times
	now := time.Now()
	exp := now.Add(time.Duration(cfg.AccessLifetime) * time.Second)

	//Create an access token object
	atoken := token.NewToken(
		usr.ID,
		env.ID,
		token.TokenTypeACCESS,
		exp,
		&now,
	)

	//Encrypt the access token and generate the cookie strings
	domain := "" //TODO: add this field eventually
	ate, atCookie := atoken.CryptAndCookie(
		env.SK,
		"/",
		domain,
		persistent,
	)
	atECookie := atoken.ExprCookie(
		"/", domain,
		cfg.ExprMultiplier, persistent,
	)
} */
