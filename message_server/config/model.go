package config

import (
	"fmt"
	"os"
	"path/filepath"
)

//
//-- Interface: IConfig
//

// Generic interface for configuration models.
type IConfig interface{}

/*
Initializes a configuration object via reading one that exists or creating
a new one with the default values.
*/
func initHelper[T IConfig](cfg *T, path string, marshaller func(*T) ([]byte, error), unmarshaller func([]byte, *T) error) error {
	//Check if the config is in a subdirectory
	ppath := filepath.Dir(path)
	if ppath != "." {
		//Create the directory structure if the path is nonexistent
		pinfo, _ := os.Stat(ppath)
		if pinfo == nil {
			err := os.MkdirAll(ppath, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	//Ensure the path is not a directory
	finfo, _ := os.Stat(path)
	if finfo != nil && finfo.IsDir() {
		return fmt.Errorf("input path '%s' points to a directory, not a file", path)
	}

	//Attempt to load the config at the given path
	cfgf, oerror := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if oerror != nil {
		return oerror
	}
	defer cfgf.Close()

	//Check if the config was previously nonexistent
	if finfo == nil {
		//Marshal the default version of the config to a byte array
		data, err := marshaller(cfg)

		//Ensure there are no errors and write the data to the file
		if err != nil {
			return err
		}
		cfgf.Write(data)

		//Return the default config object
		return nil
	}

	//Read the configuration file
	cbytes := make([]byte, finfo.Size())
	_, rerr := cfgf.Read(cbytes)
	if rerr != nil {
		return rerr
	}

	//Unmarshal a configuration object from the file data
	if err := unmarshaller(cbytes, cfg); err != nil {
		return err
	}

	//No errors, so return nil
	return nil
}
