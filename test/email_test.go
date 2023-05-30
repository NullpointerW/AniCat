package test

import (
	// "crypto/tls"

	// "crypto/tls"

	// "errors"
	// "fmt"
	// "io"
	// "net/http"
	// "net/smtp"
	// "net/url"

	"fmt"
	"testing"

	CFG "github.com/NullpointerW/mikanani/conf"

	P "github.com/NullpointerW/mikanani/pusher"
	"github.com/NullpointerW/mikanani/pusher/email"
	"github.com/NullpointerW/mikanani/pusher/wechat"
	// "gopkg.in/gomail.v2"
)

func TestEmailPush(t *testing.T) {
	var pusher P.Pusher
	pusher = email.Sender{} 
	fmt.Println(CFG.Env.Pusher)
	err := pusher.Push(P.Payload{
		114514,
		"test",
		"test.mp4",
		1145141112,
	})
	if err != nil {
		t.Error(err)
	}
}
func TestWechat(t *testing.T){
	wechat.Send()
}

