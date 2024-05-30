package tests

import (
	"fmt"
	"testing"

	"wraith.me/message_server/config"
)

const ENV_PATH = "./secret.env"

func TestNewEnv(t *testing.T) {
	env, err := config.EnvInit(ENV_PATH)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("env: %v\n", env)

	//os.Remove(ENV_PATH)
}
