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

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/pusher"

	P "github.com/NullpointerW/anicat/pusher"
	"github.com/NullpointerW/anicat/pusher/email"
	"github.com/NullpointerW/anicat/pusher/wechat"
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
		"S01E02",
		222,
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
