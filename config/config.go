package config

import "os"

type configuration struct {
	slack struct {
		token         string
		signingSecret string
	}
}

var conf = &configuration{}

func init() {
	conf.slack.token = os.Getenv("SLACK_BOT_TOKEN")
	conf.slack.signingSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

func SlackToken() string {
	return conf.slack.token
}

func SlackSigningSecret() string {
	return conf.slack.signingSecret
}
