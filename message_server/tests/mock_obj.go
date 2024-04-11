package tests

import (
	"slices"
	"time"

	"wraith.me/message_server/db/mongoutil"
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
var foo1 = Foo{
	ID:            mongoutil.MustNewUUID4(),
	Name:          "John Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"carrots", "apples", "pasta"},
}
var foo2 = Foo{
	ID:            mongoutil.MustNewUUID4(),
	Name:          "Jane Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"bananas", "melons", "ice-cream"},
}
var foo3 = Foo{
	ID:            mongoutil.MustNewUUID4(),
	Name:          "Jin Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"ramen", "rice", "sushi"},
}
