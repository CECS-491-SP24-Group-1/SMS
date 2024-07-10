package config

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	ccrypto "wraith.me/message_server/crypto"
	"wraith.me/message_server/util"
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
	ID util.UUID `env:"ID"`

	//The server's private cryptographic key.
	SK ccrypto.Privkey `env:"SK"`
}

// Overrides the `defaultPathName()` method in `IConfig`.
func (Env) defaultPathName() string {
	return DEFAULT_ENV_PATH
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

			//Special case: encode byte arrays as base64
			var vstr string
			if field.Type == reflect.TypeOf(c.SK) {
				//Get a byte slice of the crypto key
				slice := reflect.ValueOf(*c).Field(i).Interface().(ccrypto.Privkey)

				//Convert the slice to base64
				vstr = base64.RawStdEncoding.EncodeToString(slice[:])
			} else {
				//Get the value of the key as a string
				vstr = fmt.Sprintf("%v", reflect.ValueOf(*c).Field(i))
			}

			//Add the k/v pair to the output string
			fmt.Fprintf(&out, "%s=%s\n", keyn, vstr)
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

		//Decode the private key
		sks, err := base64.RawStdEncoding.DecodeString(em["SK"])
		if err != nil {
			return err
		}

		//Set the struct fields from the map and return no error
		*c = Env{
			ID: util.UUIDFromString(em["ID"]),
			SK: ccrypto.Privkey(sks),
		}
		return nil
	}

	//Create a new blank env object and set defaults
	cfg := Env{}
	cfg.ID = util.MustNewUUID7()
	var err error = nil
	_, cfg.SK, err = ccrypto.NewKeypair(nil) //`PrivateKey`` contains the public key already
	if err != nil {
		panic(err)
	}

	//Call the helper and return the results
	err = initHelper[Env](&cfg, path, marshaller, unmarshaller)
	return cfg, err
}
