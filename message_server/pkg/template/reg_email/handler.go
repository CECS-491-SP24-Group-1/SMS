package reg_email

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
	"wraith.me/message_server/pkg/config"
	"wraith.me/message_server/pkg/email"
	"wraith.me/message_server/pkg/schema/user"
	"wraith.me/message_server/pkg/util"
)

var (
	//go:embed template.html
	htmlFile     []byte
	htmlTemplate *template.Template = template.Must(
		template.New("html").Parse(string(htmlFile)),
	)
)

//
//-- CLASS: RETemplate
//

// Defines the fillable fields of the post-registration email template.
type Template struct {
	//Public fields (initial)
	ID            string
	UName         string
	Email         string
	PKFingerprint string
	PurgeTime     string

	//Public fields (filled later)
	ChallengeLink string
	ClientBaseUrl string

	//Private fields
	cfg config.Config
}

//-- Constructors

// Creates a new registration email from a user object.
func NewRegEmail(user user.User, tzOff int, chall string, cfg config.Config) Template {
	//Compose the initial template
	out := Template{
		//Initialize public fields
		ID:            user.ID.String(),
		UName:         user.Username,
		Email:         user.Email,
		PKFingerprint: user.Pubkey.Fingerprint(),
		PurgeTime:     util.Time2Offset(user.Flags.PurgeBy, tzOff).Format(time.RFC1123Z),

		//Initialize private fields
		cfg: cfg,
	}

	//Configure additional fields
	out.ChallengeLink = fmt.Sprintf(
		"%s/challenges/email/%s",
		cfg.Server.BaseUrl,
		chall,
	)
	out.ClientBaseUrl = cfg.Client.BaseUrl

	//Return the full template
	return out
}

//-- Methods

// Composes and sends an email to the email address given in the template.
func (t Template) Send() error {
	//Compose a new email
	emsg := mail.NewMSG()
	emsg.SetFrom(t.cfg.Email.Username)
	emsg.AddTo(t.Email)
	emsg.SetSubject("Your Wraith Account")

	//Create the body of the email from the template and add it to the email
	ebody, err := t.Generate()
	if err != nil {
		return err
	}
	emsg.SetBody(mail.TextHTML, ebody)

	//Send the email
	if emsg.Error != nil {
		return emsg.Error
	}
	if err := email.GetInstance().SendEmail(emsg); err != nil {
		return err
	}
	return nil
}

// Generates the body of the email using the HTML template.
func (t Template) Generate() (string, error) {
	var ebody bytes.Buffer
	if err := htmlTemplate.Execute(&ebody, t); err != nil {
		return "", err
	}
	return ebody.String(), nil
}
