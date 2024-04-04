package registration_email

//
//-- CLASS: RETemplate
//

// Defines the fillable fields of the post-registration email template.
type Template struct {
	UUID          string
	UName         string
	Email         string
	PKFingerprint string
	PurgeTime     string
	ChallengeLink string
}
