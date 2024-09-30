package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"wraith.me/message_server/db"
	"wraith.me/message_server/email"
	"wraith.me/message_server/obj/token"
	"wraith.me/message_server/redis"
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
		BindAddr   string `toml:"bind_addr" env:"SRV_BIND_ADDR" default:"127.0.0.1"`
		ListenPort int    `toml:"listen_port" env:"SRV_LISTEN_PORT" default:"8888"`
		BaseUrl    string `toml:"base_url" env:"SRV_BASE_URL" default:"http://127.0.0.1:8888/api"`
	} `toml:"server"`

	//Client configuration
	Client struct {
		BaseUrl string `toml:"base_url" env:"CLI_BASE_URL" default:"http://127.0.0.1:8080"`
	} `toml:"client"`

	//Logging configuration
	//Logging Logging `toml:"logging"`

	//Access logging configuration

	//MongoDB configuration
	MongoDB db.MConfig `toml:"mongo_db"`

	//Redis configuration
	Redis redis.RConfig `toml:"redis"`

	//SMTP configuration
	Email email.EConfig `toml:"email"`

	//Token configuration
	Token token.TConfig `toml:"token"`
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
	//Create a new blank config object and set defaults
	cfg := Config{}
	defaults.Set(&cfg)

	//Get the default path is one wasn't specified
	if path == "" {
		path = cfg.defaultPathName()
	}

	//Check if a config file exists
	exists, err := func(path string) (bool, error) {
		//Ensure the path exists
		//After this point, it is assumed that the file exists
		finfo, err := os.Stat(path)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		//The file exists, but it's a directory
		if finfo != nil && finfo.IsDir() {
			return true, fmt.Errorf("input path '%s' points to a directory, not a file", path)
		}

		//The file exists and is a valid file; further errors may be thrown later when reading
		return true, nil
	}(path)
	if err != nil {
		return cfg, err
	}

	//Create the config reader and feeders
	tomlFeeder := feeder.Toml{Path: path}
	envFeeder := feeder.Env{}
	cr := config.New().AddStruct(&cfg)

	//Add the TOML feeder only if the config exists
	if exists {
		cr.AddFeeder(tomlFeeder)
	}
	cr.AddFeeder(envFeeder)

	//Feed the config reader
	if err := cr.Feed(); err != nil {
		return cfg, err
	}

	//If the config didn't previously exist, then save the TOML file to disk
	if !exists {
		//Create the new config file
		file, err := os.Create(path)
		if err != nil {
			return cfg, err
		}
		defer file.Close()

		//Marshal the struct to TOML and write it to the config file
		encoder := toml.NewEncoder(file)
		encoder.Indent = "\t" //Set tab as the indent character
		if err := encoder.Encode(&cfg); err != nil {
			return cfg, err
		}
	}

	//Return the populated config object
	return cfg, nil
}
