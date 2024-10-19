package tests

import (
	"slices"
	"time"

	"wraith.me/message_server/pkg/util"
)

// Define the object
type Foo struct {
	ID            util.UUID `bson:"_id" json:"_id"`
	Name          string    `bson:"name" json:"name"`
	Birthday      time.Time `bson:"birthday" json:"birthday"`
	FavoriteFoods []string  `bson:"favorite_foods" json:"favorite_foods"`
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
		ID:            util.MustNewUUID4(),
		Name:          "John Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"carrots", "apples", "pasta"},
	}
	foo2 = Foo{
		ID:            util.MustNewUUID4(),
		Name:          "Jane Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"bananas", "melons", "ice-cream"},
	}
	foo3 = Foo{
		ID:            util.MustNewUUID4(),
		Name:          "Jin Doe",
		Birthday:      time.Now().Round(0),
		FavoriteFoods: []string{"ramen", "rice", "sushi"},
	}
)
