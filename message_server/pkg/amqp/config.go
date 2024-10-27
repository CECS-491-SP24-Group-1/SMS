package amqp

import (
	"fmt"

	"github.com/creasty/defaults"
)

// Configuration object for the AMQP client.
type AMQPConfig struct {
	//The address that the AMQP server is located at.
	Host string `toml:"url" env:"AMQP_HOST" default:"127.0.0.1"`

	//The port that the AMQP server is listening on.
	Port int `toml:"port" env:"AMQP_PORT" default:"5672"`

	//The username to use when connecting.
	Username string `toml:"username" env:"AMQP_USERNAME" default:""`

	//THe password to use when connecting.
	Password string `toml:"password" env:"AMQP_PASSWORD" default:""`
}

func DefaultAMQPConfig() *AMQPConfig {
	obj := &AMQPConfig{}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}
	return obj
}

// Constructs the connection URL from the config object.
func (cfg AMQPConfig) ConnURL() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
}
