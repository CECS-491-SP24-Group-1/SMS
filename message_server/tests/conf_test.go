package tests

import (
	"fmt"
	"testing"

	"wraith.me/message_server/config"
)

func TestCfgFromEnv(t *testing.T) {
	//Set some env vars
	t.Setenv("SRV_BIND_ADDR", "192.168.0.100")
	t.Setenv("RED_PASSWORD", "password")
	t.Setenv("TOK_DOMAIN", "wraithapp.me")

	//Try to load a config object
	path := "config.toml"
	cfg, err := config.ConfigInit(path)
	if err != nil {
		t.Fatal(err)
	}

	//Print the config
	fmt.Printf("%+v\n", cfg)

	//Cleanup
	//os.Remove(path)
}
