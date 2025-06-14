package redis

import "github.com/creasty/defaults"

/*
Configuration object for Redis. See the following documentation
page for more info on the available options.
https://redis.uptrace.dev/guide/go-redis.html#connecting-to-redis-server
*/
type RConfig struct {
	//The address that the Redis server is located at.
	Host string `toml:"host" env:"RED_HOST" default:"127.0.0.1"`

	//The port that the Redis server is listening on.
	Port int `toml:"port" env:"RED_PORT" default:"6379"`

	//The database that Redis should use for its KV store.
	DB int `toml:"db" env:"RED_DB" default:"0"`

	//The username to connect to the database with.
	Username string `toml:"username" env:"RED_USERNAME" default:""`

	//The password to connect to the database with.
	Password string `toml:"password" env:"RED_PASSWORD" default:""`

	//The default expiration time of entries inserted into the database.
	Expiry int64 `toml:"expiry" env:"RED_EXPIRY" default:"0"`
}

func DefaultRConfig() *RConfig {
	obj := &RConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
