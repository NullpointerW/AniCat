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
	"github.com/NullpointerW/mikanani/pusher"

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
		8964,
		"凉宫春日的忧郁",
		"[ANi] 江戶前精靈 - 04 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4",
		70891200017,
	})
	if err != nil {
		t.Error(err)
	}
}
func TestWechat(t *testing.T) {
	var pusher pusher.Pusher = wechat.WechatSender{
		Token: "xxxx",
		Uid:   "test",
	}
	pusher.Push(P.Payload{})
}
