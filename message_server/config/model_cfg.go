package config

import (
	"strings"

	"github.com/creasty/defaults"
	"github.com/pelletier/go-toml"
)

//
//-- CLASS: Config
//

// The default path at which the configuration is expected to reside.
const DEFAULT_TCONF_PATH = "./config.toml"

// Defines the configuration model for the entire application.
type Config struct {
	//Config implements the IConfig interface
	IConfig

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

// Overrides the `defaultPathName()` method in `IConfig`.
func (Config) defaultPathName() string {
	return DEFAULT_TCONF_PATH
}

// Server config block

// Logging config block
/*
type Logging struct {
	//AccessLog bool `toml:"access_log" default:"true"`
	//LogLevel
}
*/

// Configures a new TOML config object.
func ConfigInit(path string) (Config, error) {
	//Define the marshalling and unmarshalling functions
	marshaller := func(c *Config) ([]byte, error) {
		dcfgToml, err := toml.Marshal(c) //TODO: Maybe replace with toml encoder api
		if err != nil {
			return nil, err
		}

		//Format the toml string and return the results
		tomlStr := strings.TrimSpace(string(dcfgToml))
		tomlStr = strings.Replace(tomlStr, "  ", "\t", -1)
		return []byte(tomlStr), nil
	}
	unmarshaller := func(b []byte, c *Config) error {
		return toml.Unmarshal(b, c)
	}

	//Create a new blank config object and set defaults
	cfg := Config{}
	defaults.Set(&cfg)

	//Call the helper and return the results
	err := initHelper[Config](&cfg, path, marshaller, unmarshaller)
	return cfg, err
}
