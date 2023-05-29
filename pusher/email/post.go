package email

import (
	"crypto/tls"

	CFG "github.com/NullpointerW/mikanani/conf"
	"github.com/NullpointerW/mikanani/pusher"
	"gopkg.in/gomail.v2"
)

var (
	sender   gomail.Dialer
	from, to string
)

type e struct{}

func (_ e) Push(p pusher.Payload) error {
	m := gomail.NewMessage()
	// Sender
	m.SetHeader("From", from)
	// Recipient(s), can be multiple recipients, but must use the same SMTP connection.
	m.SetHeader("To", to)

	m.SetHeader("Subject", "[Mikan]剧集推送更新提醒")

	// The meaning of text/html is to set the content-type of the file as text/html,
	// and the browser will automatically call the HTML parser to process the file accordingly when it is obtained.
	// Text formatting can be specially processed using text/html, such as line breaks, indentation, bolding, etc.
	m.SetBody("text/html", Parse(p))

	if err := sender.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func init() {
	sender := gomail.NewDialer(
		CFG.Env.Pusher.Email.Host,
		CFG.Env.Pusher.Email.Port,
		CFG.Env.Pusher.Email.Username,
		CFG.Env.Pusher.Email.Password,
	)
	from, to = CFG.Env.Pusher.Email.Username, CFG.Env.Pusher.Email.Username
	if CFG.Env.Pusher.Email.SkipSSL {
		sender.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
}
