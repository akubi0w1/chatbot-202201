package slack

import "github.com/akubi0w1/chatbot-202201/code"

func ValidateContactRequest(v map[string]interface{}) error {
	if v["type"] == "" {
		return code.Errorf(code.InvalidArgument, "invalid type: type is required")
	}
	if v["content"] == "" {
		return code.Error(code.InvalidArgument, "invalid content: content is required")
	}
	return nil
}
