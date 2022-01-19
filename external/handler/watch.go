package handler

import (
	"time"

	cslack "github.com/akubi0w1/chatbot-202201/external/slack"
	"github.com/akubi0w1/chatbot-202201/log"
	"github.com/slack-go/slack"
)

type Watch struct {
	log   *log.Logger
	slack *cslack.SlackClient
}

func NewWatch(s *cslack.SlackClient) *Watch {
	return &Watch{
		log:   log.New(),
		slack: s,
	}
}

func (h *Watch) Handle(date time.Time) error {
	now := time.Now().Local()

	switch now.Format("15:04") {
	case "09:59", "10:00", "10:01":
		_, err := h.slack.PostPublicMessage(cslack.Ch_Admin, slack.MsgOptionAttachments(slack.Attachment{
			Title: "おはようございます",
			Text:  "本日は開発朝会があります！ミーティングルームにお集まりください！",
		}))
		if err != nil {
			h.slack.PostError(err)
			return err
		}
	}

	return nil
}
