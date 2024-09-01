//https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go

package tests

import (
	"fmt"
	"testing"

	"wraith.me/message_server/schema/user"
	"wraith.me/message_server/util/ms"
)

// Generally better than MS.
func BenchmarkRedactJsonDM(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if err := redactJson(ms.RedactJsonDM, true); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedactJsonMS(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if err := redactJson(ms.RedactJsonMS, true); err != nil {
			b.Fatal(err)
		}
	}
}

func TestRedactJsonDM(t *testing.T) {
	if err := redactJson(ms.RedactJsonDM, false); err != nil {
		t.Fatal(err)
	}
}

func TestRedactJsonMS(t *testing.T) {
	if err := redactJson(ms.RedactJsonMS, false); err != nil {
		t.Fatal(err)
	}
}

// Backend function for tests and benchmarks
func redactJson(redactor func(target user.User, whitelist bool, fieldNames ...string) ([]byte, error), silent bool) error {
	//Get a random user
	usr, err := GetRandomUser()
	if err != nil {
		return err
	}

	//Marshal to JSON and redact fields (blacklist)
	bjson, err := redactor(*usr, false, "options", "flags", "tokens", "email", "last_ip", "last_login")
	if err != nil {
		return nil
	}

	//Marshal to JSON and redact fields (whitelist)
	wjson, err := redactor(*usr, true, "flags", "id")
	if err != nil {
		return err
	}

	//Print results
	if !silent {
		fmt.Println(string(bjson))
		fmt.Println()
		fmt.Println(string(wjson))
	}
	return nil
}
