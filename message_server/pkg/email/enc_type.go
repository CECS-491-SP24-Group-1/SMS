//go:generate go-enum --marshal --forceupper --mustparse --nocomments --names --values
package email

//
//-- ENUM: EncType
//

/*
Defines the encryption type to use for emails. Refer to the following
file for the full list of values: https://pkg.go.dev/github.com/xhit/go-simple-mail/v2#Encryption
*/
/*
ENUM(
	NONE = 0		//No encryption will be applied. This is not recommended.
	REQUIRE = 3		//Encryption will be mandated. This should only be used if you care more about security than delivery success.
	STARTTLS = 4	//Opportunistic encryption will be applied. This is the recommended option.
)
*/
type EncType int8
