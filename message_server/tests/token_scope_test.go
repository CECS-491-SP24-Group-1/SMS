package tests

import (
	"math/rand"
	"slices"
	"testing"

	"wraith.me/message_server/obj"
)

func TestTokenScopeIsMasked(t *testing.T) {
	//Run x tests
	tests := 10
	for i := 0; i < tests; i++ {
		//Get a random number and the expected result
		rand := obj.TokenScope(rand.Intn(256))
		expected := !slices.Contains(obj.TokenScopeValues(), rand)
		actual := rand.IsMasked()

		//Test for correctness
		if expected != actual {
			t.Fatalf("[Test #%d/%d]: Incorrect mask status %v, expected %v (val: %d)\n", i+1, tests, actual, expected, rand)
		}
	}
}

func TestTokenScopeSet(t *testing.T) {
	//Test masking with dupes
	items := []obj.TokenScope{obj.TokenScopeUSER, obj.TokenScopePOSTSIGNUP, obj.TokenScopeUSER, obj.TokenScopePOSTSIGNUP, obj.TokenScopePOSTSIGNUP, obj.TokenScopePOSTSIGNUP}
	expected := obj.TokenScopePOSTSIGNUP | obj.TokenScopeUSER
	actual := obj.TokenScopeNONE
	actual.Set(items...)

	//Test for correctness; the duped values should not change the mask
	if actual != expected {
		t.Fatalf("Incorrect mask value %d, expected %d\n", actual, expected)
	}
}

func TestTokenScopeSetAN(t *testing.T) {
	//Start with a fresh token scope
	scope := obj.TokenScopeNONE

	//Mask all
	scope.SetAll()
	expected1 := obj.TokenScopeEVERWHERE
	if actual := scope; actual != expected1 {
		t.Fatalf("Unexpected token value %s, expected %s\n", actual, expected1)
	}

	//Mask none
	scope.SetNone()
	expected2 := obj.TokenScopeNONE
	if actual := scope; actual != expected2 {
		t.Fatalf("Unexpected token value %s, expected %s\n", actual, expected2)
	}
}

func TestTokenScopeUnmask(t *testing.T) {
	//Create the starting mask
	masked := obj.CreateMaskedTokenScope(obj.TokenScopePOSTSIGNUP, obj.TokenScopeUSER)
	expected := []obj.TokenScope{obj.TokenScopePOSTSIGNUP, obj.TokenScopeUSER}
	actual := masked.Unmask()

	//Test for correctness; the duped values should not change the mask
	if !slices.Equal(actual, expected) {
		t.Fatalf("Incorrect mask value %+v, expected %+v\n", actual, expected)
	}
}

func TestTokenScopeTestFor(t *testing.T) {
	scope := obj.TokenScopePOSTSIGNUP
	ts := obj.CreateMaskedTokenScope(scope)

	if actual := ts.TestFor(scope); actual != true {
		t.Fatalf("scope not found in token; got %v expected %v", actual, true)
	}
}

func TestTokenScopeTestForAll(t *testing.T) {
	scopes := []obj.TokenScope{obj.TokenScopeUSER, obj.TokenScopePOSTSIGNUP, obj.TokenScopeUSER}
	ts := obj.CreateMaskedTokenScope(scopes...)

	if actual := ts.TestForAll(scopes...); actual != true {
		t.Fatalf("scopes not found in token; got %v expected %v", actual, true)
	}
}

func TestTokenScopeTestForAny(t *testing.T) {
	scopes := []obj.TokenScope{obj.TokenScopeUSER, obj.TokenScopeUSER}
	ts := obj.CreateMaskedTokenScope(scopes...)

	if actual := ts.TestForAny(scopes...); actual != true {
		t.Fatalf("scopes not found in token; got %v expected %v", actual, true)
	}
}

func TestTokenScopeToggle(t *testing.T) {
	//Setup
	scopes := []obj.TokenScope{obj.TokenScopeUSER, obj.TokenScopePOSTSIGNUP}
	ts := obj.CreateMaskedTokenScope(scopes...)

	//Toggle users off
	ts.Toggle(obj.TokenScopeUSER)
	expected1 := false
	if actual := ts.TestFor(obj.TokenScopeUSER); actual != expected1 {
		t.Fatalf("Incorrect status on toggled mask; got %v, expected %v", actual, expected1)
	}

	//Toggle users on
	ts.Toggle(obj.TokenScopeUSER)
	expected2 := true
	if actual := ts.TestFor(obj.TokenScopeUSER); actual != expected2 {
		t.Fatalf("Incorrect status on toggled mask; got %v, expected %v", actual, expected2)
	}

	//Multiple toggles; odd number will have a different end from start, even will have the same start and end
	toggles := []obj.TokenScope{obj.TokenScopePOSTSIGNUP, obj.TokenScopePOSTSIGNUP, obj.TokenScopePOSTSIGNUP}
	expected3 := len(toggles)%2 == 0
	ts.Toggle(toggles...)
	if actual := ts.TestFor(obj.TokenScopePOSTSIGNUP); actual != expected3 {
		t.Fatalf("Incorrect status on toggled mask; got %v, expected %v", actual, expected3)
	}
}

func TestTokenScopeUnset(t *testing.T) {
	//Setup
	scopes := []obj.TokenScope{obj.TokenScopeUSER, obj.TokenScopePOSTSIGNUP}
	ts := obj.CreateMaskedTokenScope(scopes...)

	//Set users off
	ts.Unset(obj.TokenScopeUSER)
	expected1 := false
	if actual := ts.TestFor(obj.TokenScopeUSER); actual != expected1 {
		t.Fatalf("Incorrect status on unset mask; got %v, expected %v", actual, expected1)
	}

	//Set users off again; should stay off
	ts.Unset(obj.TokenScopeUSER)
	expected2 := false
	if actual := ts.TestFor(obj.TokenScopeUSER); actual != expected2 {
		t.Fatalf("Incorrect status on unset mask; got %v, expected %v", actual, expected2)
	}
}
