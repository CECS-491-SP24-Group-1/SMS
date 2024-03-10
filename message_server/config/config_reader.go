package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/creasty/defaults"
	"github.com/pelletier/go-toml"
)

const DEFAULT_CFG_NAME = "./config.toml"

func Init(path string) (*Config, error) {
	//Defaults
	if path == "" {
		path = DEFAULT_CFG_NAME
	}

	//Check if the config is in a subdirectory
	ppath := filepath.Dir(path)
	if ppath != "." {
		//Create the directory structure if the path is nonexistent
		pinfo, _ := os.Stat(ppath)
		if pinfo == nil {
			err := os.MkdirAll(ppath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	//Ensure the path is not a directory
	finfo, _ := os.Stat(path)
	if finfo != nil && finfo.IsDir() {
		return nil, fmt.Errorf("input path '%s' points to a directory, not a file", path)
	}

	//Attempt to load the config at the given path
	cfgf, oerror := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if oerror != nil {
		return nil, oerror
	}

	//Check if the config was previously nonexistent
	if finfo == nil {
		//Serialize a default config struct to toml
		var dcfg *Config = new(Config)
		defaults.Set(dcfg)
		dcfgToml, err := toml.Marshal(dcfg) //TODO: Maybe replace with toml encoder api

		//Ensure there are no errors and write the toml data to the file
		if err != nil {
			return nil, err
		}

		//Format the toml string and write it to the output file
		var tomlStr string = strings.TrimSpace(string(dcfgToml))
		tomlStr = strings.Replace(tomlStr, "  ", "\t", -1)
		cfgf.Write([]byte(tomlStr))

		//Return the default config object
		return dcfg, nil
	}

	//Read the toml file
	var tdoc []byte = make([]byte, finfo.Size())
	_, rerr := cfgf.Read(tdoc)
	if rerr != nil {
		return nil, rerr
	}

	//Deserialize the toml file to a config object
	var cfg Config
	var dserr = toml.Unmarshal(tdoc, &cfg)
	if dserr != nil {
		return nil, dserr
	}

	//Cleanup and return the config object
	cfgf.Close()
	return &cfg, nil
}
