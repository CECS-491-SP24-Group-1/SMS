package db

import "github.com/creasty/defaults"

/*
Configuration object for MongoDB. See the following documentation
page for more info on the available options.
https://www.mongodb.com/docs/drivers/node/current/fundamentals/connection/connection-options/
*/
type MConfig struct {
	//The address that the MDB server is located at.
	Host string `toml:"host" default:"127.0.0.1"`

	//The port that the MDB server is listening on.
	Port int `toml:"port" default:"27017"`

	//The name that this server should identify itself as when connecting to the database.
	AppName string `toml:"app_name" default:"WraithAPI"`

	//The username to connect to the database with.
	Username string `toml:"username" default:""`

	//The password to connect to the database with.
	Password string `toml:"password" default:""`
}

func DefaultMConfig() *MConfig {
	obj := &MConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
