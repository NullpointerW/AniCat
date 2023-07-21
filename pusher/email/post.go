package email

import (
	"crypto/tls"
	"log"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/pusher"
	util "github.com/NullpointerW/anicat/utils"
	"gopkg.in/gomail.v2"
)

var (
	sender   *gomail.Dialer
	from, to string
)

type Sender struct{}

func (_ Sender) Push(p pusher.Payload) error {
	if sender == nil {
		util.Debugln("email push disable")
		return nil
	}
	m := gomail.NewMessage()
	// Sender
	m.SetHeader("From", from)
	// Recipient(s), can be multiple recipients, but must use the same SMTP connection.
	m.SetHeader("To", to)

	m.SetHeader("Subject", "[Anicat] 剧集推送更新提醒")

	// The meaning of text/html is to set the content-type of the file as text/html,
	// and the browser will automatically call the HTML parser to process the file accordingly when it is obtained.
	// Text formatting can be specially processed using text/html, such as line breaks, indentation, bolding, etc.
	m.SetBody("text/html", Parse(p))

	if err := sender.DialAndSend(m); err != nil {
		return errs.Custom("send email error:%w", err)
	}
	return nil
}

func init() {
	if CFG.Env.Pusher.Email.Host == "" {
		log.Println("email push disable")
		return
	}
	sender = gomail.NewDialer(
		CFG.Env.Pusher.Email.Host,
		CFG.Env.Pusher.Email.Port,
		CFG.Env.Pusher.Email.Username,
		CFG.Env.Pusher.Email.Password,
	)
	from, to = CFG.Env.Pusher.Email.Username, CFG.Env.Pusher.Email.Username
	if CFG.Env.Pusher.Email.SkipSSL {
		sender.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	log.Println("email dialer init completed")
	CFG.Env.EmailPrint()
	tmpInit()
}
