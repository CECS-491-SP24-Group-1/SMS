package tests

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	"wraith.me/message_server/template/registration_email"
)

func TestHtmlTemplate(t *testing.T) {
	//Set constants
	templateLoc := "../template/registration_email/template.html"

	//Read in the template file
	tmpl, rerr := template.ParseFiles(templateLoc)
	if rerr != nil {
		t.Fatal(rerr)
	}

	//Define template data
	data := registration_email.Template{
		UUID:          "123",
		UName:         "JohnDoe",
		Email:         "jdoe@example.com",
		PKFingerprint: "123abc",
		PurgeTime:     "10/17/2000",
		ChallengeLink: "https://example.com",
	}

	//Execute the template with your data
	var filled bytes.Buffer
	err := tmpl.Execute(&filled, data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("``%s``", filled.String())
}
