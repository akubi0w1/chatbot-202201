package handler

import (
	cslack "github.com/akubi0w1/chatbot-202201/external/slack"
	"github.com/akubi0w1/chatbot-202201/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Event struct {
	log   *log.Logger
	slack *cslack.SlackClient
}

func NewEvent(s *cslack.SlackClient) *Event {
	return &Event{
		slack: s,
		log:   log.New(),
	}
}

func (h *Event) Handle(eventsAPIEvent slackevents.EventsAPIEvent) error {
	// application
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			h.slack.PostEphemeralMessage(ev.Channel, ev.User, slack.MsgOptionAttachments(cslack.BuildDefaultMenuAttachment()...))
		}
	}
	h.slack.PostPublicMessage(cslack.Ch_Admin, slack.MsgOptionText(eventsAPIEvent.Type, false))
	return nil
}
