package globals

import (
	"github.com/redis/go-redis/v9"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/email"
	cred "wraith.me/message_server/pkg/redis"
	"wraith.me/message_server/pkg/schema/user"
)

var (
	//-- MDB collections

	// Shared user collection across the entire application.
	UC *user.UserCollection

	//-- Configs

	// Shared config object across the entire application.
	Cfg *config.Config

	// Shared env object across the entire application.
	Env *config.Env

	//-- Misc

	// Shared Redis client across the entire application.
	Rcl *redis.Client

	// Shared SMTP client across the entire application.
	Smtp *email.EClient
)

// Initializes the shared globals
func Initialize(cfg *config.Config, env *config.Env) {
	//Initialize MDB collections
	UC = user.GetCollection()

	//Initialize configs
	Cfg = cfg
	Env = env

	//Initialize misc
	Rcl = cred.GetInstance().GetClient()
	Smtp = email.GetInstance()
}
