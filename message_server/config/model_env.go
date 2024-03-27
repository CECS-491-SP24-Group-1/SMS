package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"wraith.me/message_server/db/mongoutil"
)

//
//-- CLASS: Env
//

// The default path at which the env file is expected to reside.
const DEFAULT_ENV_PATH = "./secrets.env"

// Defines the configuration model for the env files.
type Env struct {
	//Env implements the IConfig interface
	IConfig

	//The ID of the server.
	ID mongoutil.UUID `env:"ID"`
}

// Configures a new env config object.
func EnvInit(path string) (Env, error) {
	//Define the marshalling and unmarshalling functions
	marshaller := func(c *Env) ([]byte, error) {
		//Create an output string
		out := strings.Builder{}

		//Loop over all fields of the struct and get the k/v pairs for marshalling
		//See: https://stackoverflow.com/a/66511341
		fields := reflect.VisibleFields(reflect.TypeOf(*c))
		for i, field := range fields {
			//Skip the first field
			if i == 0 {
				continue
			}

			//Get the tag value and determine the name of the key
			keyn := strings.Split(field.Tag.Get("env"), ",")[0]

			//Get the value of the key as a string
			val := fmt.Sprintf("%v", reflect.ValueOf(*c).Field(i))

			//Add the k/v pair to the output string
			fmt.Fprintf(&out, "%s=%s\n", keyn, val)
		}

		//Return the output string
		return []byte(out.String()), nil

	}
	unmarshaller := func(b []byte, c *Env) error {
		//Unmarshal the bytes to a map
		em, err := godotenv.UnmarshalBytes(b)
		if err != nil {
			return err
		}

		//Set the struct fields from the map and return no error
		*c = Env{
			ID: *mongoutil.UUIDFromStringOrNil(em["ID"]),
		}
		return nil
	}

	//Create a new blank env object
	cfg := Env{}
	cfg.ID = *mongoutil.MustNewUUID7()

	//Get the default path is one wasn't specified
	if path == "" {
		path = DEFAULT_ENV_PATH
	}

	//Call the helper and return the results
	err := initHelper[Env](&cfg, path, marshaller, unmarshaller)
	return cfg, err
}
