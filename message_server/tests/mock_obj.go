package tests

import (
	"slices"
	"time"

	"wraith.me/message_server/db/mongoutil"
	"wraith.me/message_server/obj"
	c "wraith.me/message_server/obj/challenge"
)

// Define the object
type Foo struct {
	ID            mongoutil.UUID `bson:"_id"`
	Name          string         `bson:"name"`
	Birthday      time.Time      `bson:"birthday"`
	FavoriteFoods []string       `bson:"favorite_foods"`
}

func fooeq(a Foo, b Foo) bool {
	return a.ID == b.ID && a.Name == b.Name && a.Birthday == b.Birthday && slices.Equal(a.FavoriteFoods, b.FavoriteFoods)
}
func fooeqa(a []Foo, b []Foo) bool {
	return slices.EqualFunc(a, b, fooeq)
}

// Create some instances
var (
	foo1 = Foo{
		ID:            mongoutil.MustNewUUID4(),
		Name:          "John Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"carrots", "apples", "pasta"},
	}
	foo2 = Foo{
		ID:            mongoutil.MustNewUUID4(),
		Name:          "Jane Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"bananas", "melons", "ice-cream"},
	}
	foo3 = Foo{
		ID:            mongoutil.MustNewUUID4(),
		Name:          "Jin Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"ramen", "rice", "sushi"},
	}
)

// Create some mock challenges
var (
	chall1 = c.NewChallengeDeterministic(
		mongoutil.UUIDFromString("9f583aab-10f1-4dbe-8388-9daa6a086cf3"),
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.UUIDFromString("9f583aab-10f1-4dbe-8388-9daa6a086cf4"), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.UUIDFromString("9f583aab-10f1-4dbe-8388-9daa6a086cf5"), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)
	chall2 = c.NewChallengeDeterministic(
		mongoutil.UUIDFromString("072b187f-e76a-4285-a36e-e363ecdc6bab"),
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.UUIDFromString("072b187f-e76a-4285-a36e-e363ecdc6bac"), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.UUIDFromString("072b187f-e76a-4285-a36e-e363ecdc6bad"), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)
	chall3 = c.NewChallengeDeterministic(
		mongoutil.UUIDFromString("7dac53f2-0371-476f-8a06-54742f12e873"),
		c.ChallengeScopeEMAIL,
		obj.Identifiable{ID: mongoutil.UUIDFromString("7dac53f2-0371-476f-8a06-54742f12e874"), Type: obj.IdTypeUSER},
		obj.Identifiable{ID: mongoutil.UUIDFromString("7dac53f2-0371-476f-8a06-54742f12e875"), Type: obj.IdTypeUSER},
		time.Now().Add(c.DEFAULT_CHALLENGE_EXPIRY),
	)
)
