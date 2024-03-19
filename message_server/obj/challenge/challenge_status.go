//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values

package challenge

//
//-- ENUM: ChallengeStatus
//

/*
Defines the status of a challenge. All challenges start with a `PENDING`
status that can either become `PASSED` or `FAILED` depending on how the
challenge was run and if all conditions set forth by it were met.
*/
/*
ENUM(
	PENDING //A challenge that has yet to be solved. This is the starting state of a challenge.
	FAILED //A Challenge that failed to be verified for whatever reason.
	PASSED //A challenge that was successfully completed.
)
*/
type ChallengeStatus int8
