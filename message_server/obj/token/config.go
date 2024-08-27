package token

import "github.com/creasty/defaults"

//TODO: add function using Redis to count refreshes

// Configuration object for user tokens.
type TConfig struct {
	//The lifetime of access tokens (in seconds). Default: 86400 (1 day).
	AccessLifetime int `toml:"access_lifetime" default:"86400"`

	//The lifetime of refresh tokens (in seconds). Default: 604800 (7 days).
	RefreshLifetime int `toml:"access_lifetime" default:"604800"`

	//The time multiplier for the lifetime of expiry tokens.
	ExprMultiplier int `toml:"expr_multiplier" default:"4"`
}

func DefaultTConfig() *TConfig {
	obj := &TConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
