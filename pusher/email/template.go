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
     <h1>剧集【${name}】更新通知</h1>
    <p>亲爱的用户，</p>
    <p></p>
    <ul>
      <li>Subject Id:${id}</li>
      <li>文件名:${dlname}</li>
      <li>大小:${size} MB</li>
    </ul>
    <p>已下载完毕</p>
    <p><a href="https://bgm.tv/subject/${id}">番剧信息</a></p>
    <img src="【图片链接】" alt="【图片描述】">
    <p>希望您能够喜欢这一最新的剧集，您也可以在我们的网站上留下您的宝贵意见和建议。</p>
    <p>谢谢！</p>
    <p>祝好！</p>
    <p><a href="https://github.com/NullpointerW/mikanani">【Mikan】</p>
  </body>
</html>`

var Def2 = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<style>
		body {
			font-family: Arial, Helvetica, sans-serif;
			font-size: 14px;
			line-height: 1.5;
			color: #24292e;
			background-color: #f6f8fa;
			padding: 20px;
		}
		h1 {
			font-size: 20px;
			margin-top: 0;
			margin-bottom: 10px;
		}
		h2 {
			font-size: 16px;
			margin-top: 20px;
			margin-bottom: 10px;
		}
		ul {
			margin-top: 0;
			margin-bottom: 10px;
			padding-left: 20px;
		}
		li {
			margin-bottom: 5px;
		}
		a {
			color: #0366d6;
			text-decoration: none;
			border-bottom: 0px solid #dfe2e5;
			padding-bottom: 0px;
			text-decoration: none;
		}
		a:hover {
			border-bottom-style: solid;
		}
		.inc {
			font-size: 12px;
			color: #586069;
			padding-top: 10px;
		}
	</style>
</head>
<body>
	<h1>${name}</h1>
	<h2>新内容已下载</h2>
	<ul>
		<li>SubjectId: ${id}</li>
		<li>文件名: ${dlname}</li>
		<li>文件大小: ${size} KB</li>
		<li><a href="https://bgm.tv/subject/${id}">在bgm.tv上查看番剧信息</a></li>
	</ul>
	<img src="http://api.bgm.tv/v0/subjects/${id}/image?type=medium">
	<p>Enjoy it,</p>
	<p>Mikan</p>
	<div class="inc">
		<p><img src="https://github.githubassets.com/images/email/global/octocat-logo.png" width="32" style=" box-sizing: border-box; ; ; ; ; ; ; ; "><a href="https://github.com/NullpointerW/mikanani"> project page</a></p>
	</div>
</body>
</html>`

var template string

func init() {
	template = Def2
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
	tmp = strings.ReplaceAll(tmp, `${id}`, strconv.Itoa(p.SubjectId))
	tmp = strings.ReplaceAll(tmp, `${name}`, p.SubjectName)
	tmp = strings.ReplaceAll(tmp, `${dlname}`, p.DownLoadName)
	tmp = strings.ReplaceAll(tmp, `${size}`, strconv.Itoa(p.Size/1024/1024))
	return tmp
}
