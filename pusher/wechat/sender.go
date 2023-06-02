// A wechat pusher implementation for example
// Visit https://github.com/wxpusher/wxpusher-client
// for more details.
package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/NullpointerW/mikanani/pusher"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

type WechatSender struct {
	Token, Uid string
}

func (wcsender WechatSender) Push(payload pusher.Payload) error {
	b, _ := json.Marshal(payload)
	p := string(b)
	msg := model.NewMessage(wcsender.Token).
		SetContentType(2).
		SetContent(p).
		AddUId(wcsender.Token)
	msgArr, err := wxpusher.SendMessage(msg)
	fmt.Println(msgArr, err)
	return err
}
