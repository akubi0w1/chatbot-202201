package slack

import (
	"strings"

	"github.com/slack-go/slack"
)

// ParseBlockActionValue blockのフォームをmapに
func ParseBlockActionValue(v map[string]map[string]slack.BlockAction) map[string]interface{} {
	result := map[string]interface{}{}
	for _, values := range v {
		for actionID, value := range values {
			switch value.Type {
			case "plain_text_input":
				result[actionID] = value.Value
			case "static_select":
				result[actionID] = value.SelectedOption.Value
			case "multi_conversations_select":
				result[actionID] = value.SelectedConversations
			case "multi_channels_select":
				result[actionID] = value.SelectedChannels
			case "conversations_select":
				result[actionID] = value.SelectedConversation
			}
		}
	}
	return result
}

func ParseContactType(ctype string) string {
	switch ctype {
	case "bug":
		return "不具合"
	case "improve":
		return "機能改善"
	case "other":
		return "その他"
	}
	return ""
}

func ParseValues(s string) map[string]interface{} {
	result := map[string]interface{}{}
	vals := strings.Split(s, ",")
	for i := range vals {
		var value interface{}
		sep := strings.SplitN(vals[i], ":", 2)
		val := strings.TrimSpace(sep[1])
		value = val
		result[strings.TrimSpace(sep[0])] = value
	}
	return result
}
