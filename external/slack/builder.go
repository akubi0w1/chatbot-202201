package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// BuildAppMentionMessage
func BuildDefaultMenuAttachment() []slack.Attachment {
	return []slack.Attachment{
		{
			Title:      "お呼びでしょうか？",
			CallbackID: CallbackSelectMenu,
			Actions: []slack.AttachmentAction{
				{
					Name: ActionSelectContact,
					Type: "button",
					Text: "意見収集",
				},
			},
		},
	}
}

// BuildContactForm
func BuildContactForm() slack.Blocks {
	// fields
	txt := slack.NewTextBlockObject("plain_text", "内容", false, false)
	input := slack.NewPlainTextInputBlockElement(
		slack.NewTextBlockObject("plain_text", "報告内容を記載してください", false, false),
		"content",
	)
	input.Multiline = true

	typeField := slack.SelectBlockElement{
		Type:        "static_select",
		ActionID:    "type",
		Placeholder: slack.NewTextBlockObject("plain_text", "報告区分", false, false),
		Options: []*slack.OptionBlockObject{
			{
				Text:  slack.NewTextBlockObject("plain_text", ParseContactType("improve"), false, false),
				Value: "improve",
			},
			{
				Text:  slack.NewTextBlockObject("plain_text", ParseContactType("bug"), false, false),
				Value: "bug",
			},
			{
				Text:  slack.NewTextBlockObject("plain_text", ParseContactType("other"), false, false),
				Value: "other",
			},
		},
	}

	// buttons cancel or submit
	submitBtnTxt := slack.NewTextBlockObject("plain_text", "送信", false, false)
	submitBtn := slack.NewButtonBlockElement(ActionSubmitContact, "", submitBtnTxt)
	submitBtn.Style = "primary"

	cancelBtnTxt := slack.NewTextBlockObject("plain_text", "キャンセル", false, false)
	cancelBtn := slack.NewButtonBlockElement(ActionCancelContact, "", cancelBtnTxt)

	return slack.Blocks{
		BlockSet: []slack.Block{
			slack.NewActionBlock("blockType", typeField),
			slack.NewInputBlock("blockContent", txt, input),
			slack.NewActionBlock("blockBtn", submitBtn, cancelBtn),
		},
	}
}

func BuildContactAttachment(_type string, content string, user string) slack.Attachment {
	color := ""
	switch _type {
	case "improve":
		color = ColorSuccess
	case "bug":
		color = ColorDanger
	}

	return slack.Attachment{
		Title: "ご意見をいただきました。",
		Color: color,
		Text:  WrapCodeBlock(fmt.Sprintf("From: %s\n\n報告区分: %s\n\n内容:\n%s", user, ParseContactType(_type), content)),
	}
}

// WrapCodeBlock
func WrapCodeBlock(msg string) string {
	return fmt.Sprintf("```\n%s\n```", msg)
}
