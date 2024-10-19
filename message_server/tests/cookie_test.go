package tests

import (
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/pkg/util"
)

func TestCookie(t *testing.T) {
	cb := util.NewCookieBuilder("sessionId", "abc123").
		WithDomain("example.com").
		WithExpiry(time.Now().Add(24 * time.Hour)).
		SetHttpOnly().
		SetSecure().
		WithPath("/").
		SetSameSiteLax()

	fmt.Println(cb.Build())
}
