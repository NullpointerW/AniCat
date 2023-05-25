package test

import (
	// "crypto/tls"
	"crypto/tls"
	"errors"
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
	<!DOCTYPE html>
	<html>
	<head>
		<title>New Episode Available - TV Show Name</title>
	</head>
	<body>
		<div style="background-color: #f2f2f2; padding: 20px;">
			<h2 style="font-family: Arial, sans-serif; color: #333333; font-size: 24px;">New Episode Available - TV Show Name</h2>
			<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Dear Subscriber,</p>
			<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">We are excited to announce that a new episode of <strong>TV Show Name</strong> is now available for download! The file size for this episode is <strong>[File Size]</strong>.</p>
			<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">To access the latest episode, please log in to your account and navigate to the "Episodes" section. From there, you can select the newest episode and begin downloading.</p>
			<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Thank you for being a loyal subscriber, and we hope you enjoy the latest installment of <strong>TV Show Name</strong>.</p>
			<p style="font-family: Arial, sans-serif; color: #333333; font-size: 16px;">Best regards,<br>[Your Company Name]</p>
		</div>
	</body>
	</html>
	`

	// QQ 邮箱：
	// SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/994 | 非SSL协议端口：25）
	// 163 邮箱：
	// SMTP 服务器地址：smtp.163.com（端口：25）
	host := "smtp.qq.com"
	port := 25
	userName := "xxx@qq.com"
	password := "xxx"

	m := gomail.NewMessage()
	m.SetHeader("From", userName) // 发件人
	// m.SetHeader("From", "alias"+"<"+userName+">") // 增加发件人别名

	m.SetHeader("To", "xxx@qq.com") // 收件人，可以多个收件人，但必须使用相同的 SMTP 连接

	m.SetHeader("Subject", "Test!") // 邮件主题

	// text/html 的意思是将文件的 content-type 设置为 text/html 的形式，浏览器在获取到这种文件时会自动调用html的解析器对文件进行相应的处理。
	// 可以通过 text/html 处理文本格式进行特殊处理，如换行、缩进、加粗等等
	m.SetBody("text/html", message)

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
