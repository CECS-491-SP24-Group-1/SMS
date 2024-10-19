package token

import "github.com/creasty/defaults"

//TODO: add function using Redis to count refreshes

// Configuration object for user tokens.
type TConfig struct {
	//The lifetime of access tokens (in seconds). Default: 86400 (1 day).
	AccessLifetime int `toml:"access_lifetime" env:"TOK_ACCESS_LIFETIME" default:"86400"`

	//The lifetime of refresh tokens (in seconds). Default: 604800 (7 days).
	RefreshLifetime int `toml:"refresh_lifetime" env:"TOK_REFRESH_LIFETIME" default:"604800"`

	//The time multiplier for the lifetime of expiry tokens.
	ExprMultiplier int `toml:"expr_multiplier" env:"TOK_EXPR_MULTIPLIER" default:"2"`

	//The domain to use for tokens.
	Domain string `toml:"domain" env:"TOK_DOMAIN" default:"localhost"`

	//The path at which access token cookies are valid.
	AccessCookiePath string `toml:"access_cookie_path" env:"TOK_ACCESS_COOKIE_PATH" default:"/api"`

	//The path at which refresh token cookies are valid.
	RefreshCookiePath string `toml:"refresh_cookie_path" env:"TOK_REFRESH_COOKIE_PATH" default:"/api/auth"`
}

func DefaultTConfig() *TConfig {
	obj := &TConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
