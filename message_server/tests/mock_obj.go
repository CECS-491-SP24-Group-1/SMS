package tests

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

// Define the object
type Foo struct {
	ID            uuid.UUID
	Name          string
	Birthday      time.Time
	FavoriteFoods []string
}

func fooeq(a Foo, b Foo) bool {
	return a.ID == b.ID && a.Name == b.Name && a.Birthday == b.Birthday && slices.Equal(a.FavoriteFoods, b.FavoriteFoods)
}
func fooeqa(a []Foo, b []Foo) bool {
	return slices.EqualFunc(a, b, fooeq)
}

// Create some instances
var foo1 = Foo{
	ID:            uuid.New(),
	Name:          "John Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"carrots", "apples", "pasta"},
}
var foo2 = Foo{
	ID:            uuid.New(),
	Name:          "Jane Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"bananas", "melons", "ice-cream"},
}
var foo3 = Foo{
	ID:            uuid.New(),
	Name:          "Jin Doe",
	Birthday:      time.Now().Round(0),
	FavoriteFoods: []string{"ramen", "rice", "sushi"},
}
