package email

import (
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/pusher"
)

const def = `<!DOCTYPE html>
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
	<h2>剧集已更新</h2>
	<ul>
		<li>SubjectId: ${id}</li>
		<li>文件名: ${dlname}</li>
		<li>文件大小: ${size} MB</li>
		<li><a href="https://bgm.tv/subject/${id}">在bgm.tv上查看番剧信息</a></li>
	</ul>
	<img src="http://api.bgm.tv/v0/subjects/${id}/image?type=medium">
	<p>已下载完成</p>
	<p>Enjoy it,</p>
	<p>Mikan</p>
	<div class="inc">
		<p><img src="https://github.githubassets.com/images/email/global/octocat-logo.png" width="32" style=" box-sizing: border-box; ; ; ; ; ; ; ; "><a href="https://github.com/NullpointerW/anicat"> project page</a></p>
	</div>
</body>
</html>`

var template string

func init() {
	template = def
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
