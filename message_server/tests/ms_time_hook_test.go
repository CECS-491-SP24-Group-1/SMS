package tests

import (
	"fmt"
	"testing"
	"time"

	"wraith.me/message_server/util/ms"
)

// Example struct with time.Time fields
type myStruct struct {
	CreatedAt time.Time `mapstructure:"created_at"`
}

func TestMSTimeHook(t *testing.T) {
	now := time.Now()
	in := myStruct{
		CreatedAt: now,
		//UpdatedAt: now.Add(1 * time.Hour),
		//Name:      "Example",
	}

	fmt.Printf("time in:  %v\n", in)

	mp := make(map[string]interface{})
	err1 := ms.MSTextMarshal(in, &mp, "mapstructure")
	if err1 != nil {
		t.Fatal(err1)
	}

	fmt.Printf("time out: %v\n", mp)

	out := myStruct{}
	err2 := ms.MSTextUnmarshal(mp, &out, "mapstructure")
	if err2 != nil {
		t.Fatal(err2)
	}

	fmt.Printf("time cnv: %v\n", out)
}
