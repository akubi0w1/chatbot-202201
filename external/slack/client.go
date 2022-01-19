package slack

import (
	"fmt"
	"log"

	"github.com/akubi0w1/chatbot-202201/code"
	"github.com/akubi0w1/chatbot-202201/config"
	"github.com/slack-go/slack"
)

type SlackClient struct {
	cli *slack.Client
}

func NewSlackClient() *SlackClient {
	return &SlackClient{
		cli: slack.New(config.SlackToken()),
	}
}

// PostPublicMessage 公開メッセージを送信する
func (s *SlackClient) PostPublicMessage(channelID string, opts ...slack.MsgOption) (timestamp string, err error) {
	_, ts, err := s.cli.PostMessage(channelID, opts...)
	if err != nil {
		return "", err
	}
	return ts, nil
}

// PostEphemeralMessage 指定したユーザにのみ表示される一時的なメッセージを送信する
func (s *SlackClient) PostEphemeralMessage(channelID, userID string, opts ...slack.MsgOption) (timestamp string, err error) {
	ts, err := s.cli.PostEphemeral(channelID, userID, opts...)
	if err != nil {
		return "", code.Errorf(code.Slack, "failed to post ephemeral message: %v", err)
	}
	return ts, nil
}

// ReplaceOriginalMessage メッセージを置き換える
func (s *SlackClient) ReplaceOriginalMessage(channelID, responseURL string, opts ...slack.MsgOption) error {
	_opts := append(opts, slack.MsgOptionReplaceOriginal(responseURL))
	if _, _, _, err := s.cli.SendMessage(channelID, _opts...); err != nil {
		return code.Errorf(code.Slack, "fialed to replace original message: %v", err)
	}
	return nil
}

// ReplyMessage 特定のメッセージにスレッドとして返信する
func (s *SlackClient) ReplyMessage(channelID, messageTs string, opts ...slack.MsgOption) error {
	_opts := append(opts, slack.MsgOptionTS(messageTs))
	if _, _, _, err := s.cli.SendMessage(channelID, _opts...); err != nil {
		return code.Errorf(code.Slack, "failed to reply message: %v", err)
	}
	return nil
}

// UpdateMessage メッセージを置き換える
func (s *SlackClient) UpdateMessage(channelID, messageTs string, opts ...slack.MsgOption) error {
	if _, _, _, err := s.cli.UpdateMessage(channelID, messageTs, opts...); err != nil {
		return code.Errorf(code.Slack, "failed to update message: %v", err)
	}
	return nil
}

// PostError 秘書メリサのお部屋と送信者に対して、エラーがおきたよって通知をする
// channelID, responseURL, userID: 申請者への連絡を飛ばす
func (s *SlackClient) PostErrorWithNotify(channelID, responseURL, userID string, _err error) error {
	log.Printf("[ERROR] %v", _err)
	// to秘書メリサ
	_, err := s.PostPublicMessage(Ch_Admin, slack.MsgOptionAttachments(slack.Attachment{
		Title: "処理に失敗しました。エラー内容の報告になります。",
		Text:  fmt.Sprintf("```\nUserID: %s\n\nError:\n%s\n```", userID, _err),
		Color: ColorDanger,
	}))
	if err != nil {
		return err
	}

	// toエラー告知
	err = s.ReplaceOriginalMessage(channelID, responseURL, slack.MsgOptionAttachments(slack.Attachment{
		Title: "処理に失敗しました。",
		Text:  "改善のための連絡をいたしましたので、しばらくお待ちいただき再度お試しください。\nご迷惑をおかけしてしまい、申し訳ございません。"}))
	if err != nil {
		return err
	}

	return nil
}

// PostError メリサのお部屋にエラーおきたよって通知する
func (s *SlackClient) PostError(_err error) error {
	log.Printf("[ERROR] %v", _err)
	// to秘書メリサ
	_, err := s.PostPublicMessage(Ch_Admin, slack.MsgOptionAttachments(slack.Attachment{
		Title: "処理に失敗しました。エラー内容の報告になります。",
		Text:  fmt.Sprintf("```\n%s\n```", _err),
		Color: ColorDanger,
	}))
	if err != nil {
		return err
	}

	return nil
}

// GetUserInfo
func (s *SlackClient) GetUserInfo(userID string) (*slack.User, error) {
	user, err := s.cli.GetUserInfo(userID)
	if err != nil {
		return nil, code.Errorf(code.Slack, "failed to get user: %v", err)
	}
	return user, nil
}

// GetConversationInfo
func (s *SlackClient) GetConversationInfo(channelID string, includeLocale bool) (*slack.Channel, error) {
	return s.cli.GetConversationInfo(channelID, includeLocale)
}

func (s *SlackClient) GetConversationHistory(params *slack.GetConversationHistoryParameters) (*slack.GetConversationHistoryResponse, error) {
	return s.cli.GetConversationHistory(params)
}

func (s *SlackClient) GetPermaLink(channelID, timestamp string) (string, error) {
	return s.cli.GetPermalink(&slack.PermalinkParameters{
		Channel: channelID,
		Ts:      timestamp,
	})
}
