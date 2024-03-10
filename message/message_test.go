package message

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenericMessage(t *testing.T) {
	//Get a new generic message object and the current time
	timeNow := time.Now()
	msg := NewGenericMessage()

	//The time of the message's creation should match the current time
	if msg.Created().Unix() != timeNow.Unix() {
		t.Errorf("mismatched message time and current time; %d::%d", msg.Created().Unix(), timeNow.Unix())
		t.FailNow()
	}
}

func TestExpiringMessage(t *testing.T) {
	//Create a list of durations
	durations := []time.Duration{time.Second, time.Minute, time.Hour}
	min := 10
	max := 100

	//Loop x times
	tests := 10
	for i := 0; i < tests; i++ {
		//Get a random duration and multiplier
		mult := rand.Intn(max-min) + min
		dur := durations[rand.Intn(len(durations))] * time.Duration(mult)

		//Create an expiring message
		expm, _ := NewExpiringMessage(dur)

		//Test to ensure the message has the proper expiry time
		expected := time.Now().Add(dur)
		actual := expm.Expiry
		if expected.Unix() != actual.Unix() {
			t.Errorf("mismatched expected time and actual time; %d::%d", expected.Unix(), actual.Unix())
			t.FailNow()
		}

		//Debug print
		fmt.Printf("exp_m #%d: %s\n", i+1, expm)

		//Randomly expire the message
		expireNow := rand.Intn(2) != 0
		if expireNow {
			//Check if this succeeded
			expm.ExpireNow()
			if !expm.IsExpired() {
				t.Errorf("message expiration failed; remaining: %s", expm.DurationToExpiry())
				t.FailNow()
			}
		}
	}
}
