// A wechat pusher implementation for example
// Visit https://github.com/wxpusher/wxpusher-client for more details
package wechat

import (
	"fmt"

	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

func Send() {
	msg := model.NewMessage("xxx").
		SetContentType(2).
		SetContent("test").
		AddUId("UID_xxx")
	msgArr, err := wxpusher.SendMessage(msg)
	fmt.Println(msgArr, err)
}
