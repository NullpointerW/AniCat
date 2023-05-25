package test

import (
	// "crypto/tls"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"testing"

	"gopkg.in/gomail.v2"


)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func TestSendQQEmail(t *testing.T) {
	message := `
    <p> Hello %s,</p>
	     $torr$
		<p style="text-indent:2em">test test test test test test test test test test test test.</p> 
		<p style="text-indent:2em">test test test test test test test test test test test test.</p>
 
		<p style="text-indent:2em">test test test test test test test test test test test test.</P>
 
		<p style="text-indent:2em">Best Wishes!</p>
	`

	// QQ 邮箱：
	// SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/994 | 非SSL协议端口：25）
	// 163 邮箱：
	// SMTP 服务器地址：smtp.163.com（端口：25）
	host := "smtp.qq.com"
	port := 25
	userName := "xxxx@qq.com"
	password := "xxx"

	m := gomail.NewMessage()
	m.SetHeader("From", userName) // 发件人
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", "xxx@qq.com") // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接

	m.SetHeader("Subject", "Test!") // 邮件主题

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	m.SetBody("text/html", fmt.Sprintf(message, "testUser"))

	// text/plain的意思是将文件设置为纯文本的形式，浏览器在获取到这种文件时并不会对其进行处理
	// m.SetBody("text/plain", "纯文本")
	// m.Attach("test.sh")   // 附件文件，可以是文件，照片，视频等等
	// m.Attach("lolcatVideo.mp4") // 视频
	// m.Attach("lolcat.jpg") // 照片

	d := gomail.NewDialer(
		host,
		port,
		userName,
		password,
	)
	// 关闭SSL协议认证
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// auth := LoginAuth(userName, password)

	// d.Auth = auth

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}



