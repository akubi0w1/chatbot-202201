package handler

import (
	"fmt"

	cslack "github.com/akubi0w1/chatbot-202201/external/slack"
	"github.com/akubi0w1/chatbot-202201/log"
	"github.com/slack-go/slack"
)

type Action struct {
	log   *log.Logger
	slack *cslack.SlackClient
}

func NewAction(s *cslack.SlackClient) *Action {
	return &Action{
		log:   log.New(),
		slack: s,
	}
}

func (h *Action) Handle(payload slack.InteractionCallback) error {
	switch payload.Type {
	case slack.InteractionTypeInteractionMessage: // attachmentのボタンや選択肢
		switch payload.CallbackID {
		case cslack.CallbackSelectMenu:
			return h.handleSelectMenu(payload)

		// contact
		case cslack.CallbackFixedContact:
			return h.handleFixedContact(payload)
		}

	case slack.InteractionTypeBlockActions: // blockのボタンや選択肢
		switch payload.ActionCallback.BlockActions[0].ActionID {

		// contact
		case cslack.ActionSubmitContact:
			return h.handleSubmitContact(payload)
		case cslack.ActionCancelContact:
			return h.handleCancelContact(payload)
		}
	}
	return nil
}

// handleSelectMenu 初期メニューの選択をハンドリング
func (h *Action) handleSelectMenu(payload slack.InteractionCallback) error {
	switch payload.ActionCallback.AttachmentActions[0].Name {
	case cslack.ActionSelectContact:
		form := cslack.BuildContactForm()
		if err := h.slack.ReplaceOriginalMessage(payload.Channel.ID, payload.ResponseURL, slack.MsgOptionBlocks(form.BlockSet...)); err != nil {
			h.slack.PostError(err)
			return err
		}
	}
	return nil
}

// handleSubmitContact 意見収集の送信をハンドリング
func (h *Action) handleSubmitContact(payload slack.InteractionCallback) error {
	values := cslack.ParseBlockActionValue(payload.BlockActionState.Values)
	if err := cslack.ValidateContactRequest(values); err != nil {
		h.slack.PostEphemeralMessage(payload.Channel.ID, payload.User.ID, slack.MsgOptionAttachments(slack.Attachment{Text: "入力値が不正です"}))
		return err
	}

	if err := h.slack.ReplaceOriginalMessage(payload.Channel.ID, payload.ResponseURL, slack.MsgOptionAttachments(
		slack.Attachment{
			Title:      "送信内容を確認してください",
			Text:       cslack.WrapCodeBlock(fmt.Sprintf("報告区分: %s\n\n内容:\n%s", cslack.ParseContactType(values["type"].(string)), values["content"])),
			MarkdownIn: []string{"text"},
			CallbackID: cslack.CallbackFixedContact,
			Actions: []slack.AttachmentAction{
				{
					Name:  cslack.ActionFixedContact,
					Text:  "確定",
					Type:  "button",
					Value: fmt.Sprintf("type:%s,content:%s", values["type"], values["content"]),
					Style: "primary",
				},
				{
					Name: cslack.ActionCancelContact,
					Text: "キャンセル",
					Type: "button",
				},
			},
		},
	)); err != nil {
		h.slack.PostError(err)
		return err
	}

	return nil
}

// handleFixedContact 意見収集の確定をハンドリング
func (h *Action) handleFixedContact(payload slack.InteractionCallback) error {
	switch payload.ActionCallback.AttachmentActions[0].Name {
	case cslack.ActionFixedContact:
		values := cslack.ParseValues(payload.ActionCallback.AttachmentActions[0].Value)
		user, err := h.slack.GetUserInfo(payload.User.ID)
		if err != nil {
			h.slack.PostError(err)
			return err
		}

		// post manager
		_, err = h.slack.PostPublicMessage(cslack.Ch_Admin, slack.MsgOptionAttachments(cslack.BuildContactAttachment(values["type"].(string), values["content"].(string), user.Profile.DisplayName)))
		if err != nil {
			h.slack.PostError(err)
			return err
		}

		// replace
		err = h.slack.ReplaceOriginalMessage(payload.Channel.ID, payload.ResponseURL, slack.MsgOptionAttachments(slack.Attachment{
			Title: "送信が完了しました。ご報告ありがとうございます。",
			Text:  cslack.WrapCodeBlock(fmt.Sprintf("報告区分: %s\n\n内容:\n%s", cslack.ParseContactType(values["type"].(string)), values["content"])),
		}))
		if err != nil {
			h.slack.PostError(err)
			return err
		}

	case cslack.ActionCancelContact:
		if err := h.slack.ReplaceOriginalMessage(payload.Channel.ID, payload.ResponseURL, slack.MsgOptionAttachments(slack.Attachment{Text: "送信をキャンセルしました"})); err != nil {
			h.slack.PostError(err)
			return err
		}
	}
	return nil
}

// handleCancelContact 意見収集のキャンセルをハンドリング
func (h *Action) handleCancelContact(payload slack.InteractionCallback) error {
	if err := h.slack.ReplaceOriginalMessage(payload.Channel.ID, payload.ResponseURL, slack.MsgOptionAttachments(slack.Attachment{Text: "送信をキャンセルしました"})); err != nil {
		h.slack.PostError(err)
		return err
	}
	return nil
}
