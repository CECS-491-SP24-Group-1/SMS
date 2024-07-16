package db

import "github.com/creasty/defaults"

/*
Configuration object for MongoDB. See the following documentation
page for more info on the available options.
https://www.mongodb.com/docs/drivers/node/current/fundamentals/connection/connection-options/
*/
type MConfig struct {
	//The connection string to use when establishing a connection to the MongoDB server.
	ConnStr string `toml:"conn_str" default:"mongodb://127.0.0.1:27017"`

	//The timeout (in seconds) to use for connections.
	Timeout int64 `toml:"timeout" default:"10"`
}

func DefaultMConfig() *MConfig {
	obj := &MConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
