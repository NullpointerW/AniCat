package email

import (
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
	<title>剧集更新 - $name$(SubjectId:$id$)</title>
</head>
<body>
	<div style="background-color: #f2f2f2; padding: 20px;">
		<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Dear Subscriber,</p>
		<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">新增内容<strong>$dlname$</strong> 已下载完成,文件大小<strong>$size$</strong>kb</p>
		<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">To access the latest episode, please log in to your account and navigate to the "Episodes" section. From there, you can select the newest episode and begin downloading.</p>
		<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Thank you for being a loyal subscriber, and we hope you enjoy the latest installment of <strong>TV Show Name</strong>.</p>
		<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Best regards,<br>[Mikan]</p>
	</div>
</body>
</html>
`

var template string

func Init() {
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
