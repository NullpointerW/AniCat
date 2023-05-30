package email

import (
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"

	CFG "github.com/NullpointerW/mikanani/conf"
	"github.com/NullpointerW/mikanani/pusher"
)


const Default = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>剧集更新通知</title>
  </head>
  <body>
    <h1>剧集【$name$】更新通知</h1>
    <p>亲爱的用户，</p>
    <p></p>
    <ul>
      <li>Subject Id:$id$</li>
      <li>文件名:【$dlname$】</li>
      <li>大小:$size$ KB</li>
    </ul>
    <p>您可以通过以下链接直接访问我们网站，观看这一最新的剧集：</p>
    <p><a href="【剧集链接】">【剧集链接】</a></p>
    <p>希望您能够喜欢这一最新的剧集，您也可以在我们的网站上留下您的宝贵意见和建议。</p>
    <p>谢谢！</p>
    <p>祝好！</p>
    <p>【发件人名称】</p>
  </body>
</html>`

var template string

func init() {
	template = Default
	if p := CFG.Env.Pusher.Email.TemplatePath; p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			log.Printf("load email template file fail:%s \n", err)
		} else {
			template = string(b)
		}
	}
}

func Parse(p pusher.Payload) string {
	tmp := template
	tmp = strings.ReplaceAll(tmp, `$id$`, strconv.Itoa(p.SubjectId))
	tmp = strings.ReplaceAll(tmp, `$name$`, p.SubjectName)
	tmp = strings.ReplaceAll(tmp, `$dlname$`, p.DownLoadName)
	tmp = strings.ReplaceAll(tmp, `$size$`, strconv.Itoa(p.Size))
	return tmp
}
