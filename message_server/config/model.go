package config

/* Defines the configuration model for the entire application. */
type Config struct {
	//Server configuration
	Server struct {
		BindAddr   string `toml:"bind_addr" default:"127.0.0.1"`
		ListenPort int    `toml:"listen_port" default:"8888"`
	} `toml:"server"`

	//Logging configuration
	//Logging Logging `toml:"logging"`

	//Access logging configuration

	//MongoDB configuration
	MongoDB struct {
		Host string `toml:"host" default:"127.0.0.1"`
		Port int    `toml:"port" default:"27017"`
	} `toml:"mongo_db"`
}

// Server config block

// Logging config block
/*
type Logging struct {
	//AccessLog bool `toml:"access_log" default:"true"`
	//LogLevel
}
*/
