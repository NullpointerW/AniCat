package wechat

import (
	"fmt"

	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

func Send() {
	msg := model.NewMessage("xx").
	    SetContentType(2).
		SetContent(`<!DOCTYPE html>
		<html>
		  <head>
			<meta charset="utf-8">
			<title>剧集更新通知</title>
		  </head>
		  <body>
			<h1>剧集更新通知</h1>
			<p>亲爱的用户，</p>
			<p>我们网站上的一部电视剧已经更新了新的剧集：</p>
			<ul>
			  <li>电视剧名称：【剧集名称】</li>
			  <li>更新集数：【更新集数】</li>
			  <li>更新日期：【更新日期】</li>
			</ul>
			<p>您可以通过以下链接直接访问我们网站，观看这一最新的剧集：</p>
			<p><a href="【剧集链接】">【剧集链接】</a></p>
			<p>以下是本剧的剧照：</p>
			<img src="【图片链接】" alt="【图片描述】">
			<p>希望您能够喜欢这一最新的剧集，您也可以在我们的网站上留下您的宝贵意见和建议。</p>
			<p>谢谢！</p>
			<p>祝好！</p>
			<p>【发件人名称】</p>
		  </body>
		</html>`).
		AddUId("x")
	msgArr, err := wxpusher.SendMessage(msg)
	fmt.Println(msgArr, err)
}
