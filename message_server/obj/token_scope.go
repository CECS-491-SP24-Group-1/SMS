//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package obj

import (
	"slices"
)

//
//-- ENUM: TokenScope
//

/*
Defines the scope for which the token is valid. These values span from
lowest permissions to highest permissions.
*/
/*
ENUM(
	NONE =			0		//The token is valid nowhere.												(binary: 0000 0000)
	POST_SIGNUP = 	1		//The token is only allowed to complete the login challenges.				(binary: 0000 0001)
	USER =			2		//The token is allowed to be used everywhere that a normal user can access	(binary: 0000 0010)
	EVERWHERE =		255 	//The token is valid anywhere.												(binary: 1111 1111)
)
*/
type TokenScope uint8

/*
Creates a masked token scope, given the scopes for which the token is
valid and the starting value. If duplicated validity scopes are defined,
then they are simply ignored, due to the nature of the bitwise OR operation.
*/
func CreateMaskedTokenScopeFI(init TokenScope, scopes ...TokenScope) TokenScope {
	//Loop over the input scopes
	for _, scope := range scopes {
		//Bitwise OR the current value with the output
		init |= scope
	}

	//Return the bitmasked value
	return init
}

/*
Creates a masked token scope, given the scopes for which the token is
valid. If duplicated validity scopes are defined, then they are simply
ignored, due to the nature of the bitwise OR operation.
*/
func CreateMaskedTokenScope(scopes ...TokenScope) TokenScope {
	return CreateMaskedTokenScopeFI(0, scopes...)
}

/*
Determines if a token scope is masked. This is determined the be the case
if the current value doesn't match any of the constants.
*/
func (ts TokenScope) IsMasked() bool {
	return !slices.Contains(TokenScopeValues(), ts)
}

/*
Masks the current token scope with a list of given scopes for which the
token is valid. If the masked token already is valid for the given scope,
then it is ignored due to how the bitwise OR operator works.
*/
func (ts *TokenScope) Set(scopes ...TokenScope) {
	*ts = CreateMaskedTokenScopeFI(*ts, scopes...)
}

// Makes the token valid for all scopes.
func (ts *TokenScope) SetAll() {
	*ts = TokenScopeEVERWHERE
}

// Makes the token valid for no scopes.
func (ts *TokenScope) SetNone() {
	*ts = TokenScopeNONE
}

// Tests to see if a token is valid for the given scope.
func (ts TokenScope) TestFor(scope TokenScope) bool {
	return ts&scope == scope
}

/*
Tests to see if all of the given scopes are valid for the scope object.
All scopes given must be included in the scope mask for this function
to return true.
*/
func (ts TokenScope) TestForAll(scopes ...TokenScope) bool {
	//Loop over all the given scopes
	for _, scope := range scopes {
		//Check if the current scope is not in the masked value
		if !ts.TestFor(scope) {
			return false
		}
	}

	//No non-included scopes so return true
	return true
}

/*
Tests to see if any of the given scopes are valid for the scope object.
At least one scope given must be included in the scope mask for this
function to return true.
*/
func (ts TokenScope) TestForAny(scopes ...TokenScope) bool {
	//Loop over all the given scopes
	for _, scope := range scopes {
		//Check if the current scope is in the masked value
		if ts.TestFor(scope) {
			return true
		}
	}

	//No matches, so return false
	return false
}

// Toggles the values of multiple scopes from the token using a bitwise XOR.
func (ts *TokenScope) Toggle(scopes ...TokenScope) {
	//Loop over the input scopes
	for _, scope := range scopes {
		//Bitwise XOR the token scope with the current scope
		*ts ^= scope
	}
}

/*
Unmasks a masked token scope, getting all the scopes for which the token
is valid, with the exception of the `NONE` and `EVERYWHERE` scopes.
*/
func (ts TokenScope) Unmask() []TokenScope {
	//Create the output array
	out := make([]TokenScope, 0)

	//Loop over all possible enum values
	for _, scope := range tsValuesExcludeAN() {
		//Check if the current scope is masked in the input token
		if ts&scope == scope {
			//Append the scope to the output list
			out = append(out, scope)
		}
	}

	//Return the unmasked array
	return out
}

// Unsets the values of multiple scopes from the token using a bitwise NAND.
func (ts *TokenScope) Unset(scopes ...TokenScope) {
	//Loop over the input scopes
	for _, scope := range scopes {
		//Bitwise NAND the token scope with the current scope
		*ts &= ^scope
	}
}

// Returns a list of all token scope values except for none and all.
func tsValuesExcludeAN() []TokenScope {
	vals := TokenScopeValues()
	vals = vals[1 : len(vals)-1] //First and last items are the values to remove
	return vals
}
