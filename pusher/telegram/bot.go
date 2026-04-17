package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	CFG "github.com/NullpointerW/anicat/conf"
	"github.com/NullpointerW/anicat/log"
	"github.com/NullpointerW/anicat/pusher"
)

var Poster *Bot

type Bot struct {
	token  string
	chatId string
	client *http.Client
}

func (b *Bot) Push(p pusher.Payload) error {
	if b == nil {
		log.Debug(nil, "telegram push disable")
		return nil
	}
	text := fmt.Sprintf(
		"[AniCat] 剧集更新提醒\n名称: %s\n剧集: %s\n文件: %s\n大小: %d MB",
		p.SubjectName, p.Episode, p.DownLoadName, p.Size,
	)
	body, _ := json.Marshal(map[string]string{
		"chat_id": b.chatId,
		"text":    text,
	})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)
	resp, err := b.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telegram push: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram push: unexpected status %d", resp.StatusCode)
	}
	log.Info(log.Struct{"name", p.SubjectName, "epi", p.Episode}, "telegram push sent")
	return nil
}

func init() {
	if CFG.SrvCTL {
		return
	}
	if CFG.Env.Pusher.Telegram.Token == "" {
		log.Warn(nil, "telegram push disable")
		return
	}
	Poster = &Bot{
		token:  CFG.Env.Pusher.Telegram.Token,
		chatId: CFG.Env.Pusher.Telegram.ChatId,
		client: &http.Client{},
	}
	log.Info(log.Struct{"chat_id", CFG.Env.Pusher.Telegram.ChatId}, "telegram bot init completed")
}

var _ pusher.Pusher = (*Bot)(nil)
