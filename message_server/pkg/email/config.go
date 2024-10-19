package email

import "github.com/creasty/defaults"

// Configuration object for the email daemon.
type EConfig struct {
	//Whether the emailer functionality is enabled.
	Enabled bool `toml:"enabled" env:"EMAIL_ENABLED" default:"true"`

	//The address that the email server is located at.
	Host string `toml:"host" env:"EMAIL_HOST" default:"127.0.0.1"`

	//The port that the email server is listening on.
	Port int `toml:"port" env:"EMAIL_PORT" default:"587"`

	//The username to connect to the email server with.
	Username string `toml:"username" env:"EMAIL_USERNAME" default:""`

	//The password to connect to the email server with.
	Password string `toml:"password" env:"EMAIL_PASSWORD" default:""`

	//The encryption type to use for the outgoing emails.
	EncType EncType `toml:"enc_type" env:"EMAIL_ENC_TYPE" default:"STARTTLS"`

	//Whether the certificate of the server should be verified. It is a good idea to not turn this off. Only toggle if you are ABSOLUTELY sure.
	VerifyCert bool `toml:"verify_cert" env:"EMAIL_VERIFY_CERT" default:"true"`
}

func DefaultEConfig() *EConfig {
	obj := &EConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}
