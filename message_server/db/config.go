package db

import "github.com/creasty/defaults"

/*
Configuration object for MongoDB. See the following documentation
page for more info on the available options.
https://www.mongodb.com/docs/drivers/node/current/fundamentals/connection/connection-options/
*/
type MConfig struct {
	//THe l
	Host     string `toml:"host" default:"127.0.0.1"`
	Port     int    `toml:"port" default:"27017"`
	AppName  string `toml:"app_name" default:"WraithAPI"`
	Username string `toml:"username" default:""`
	Password string `toml:"password" default:""`
}

func DefaultMConfig() *MConfig {
	obj := &MConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
