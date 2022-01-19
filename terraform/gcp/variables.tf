locals {
    slack_workspace = "akubi-post"
    service = "chatbot"
}

variable "project_id" {
  type    = string
  default = "hoge"
}

variable "project_number" {
  type    = string
  default = "hoge"
}

variable "region" {
  type    = string
  default = "asia-northeast1"
}

variable "zone" {
  type    = string
  default = "asia-northeast1-b"
}

variable "slack_bot_token" {
  type = string
  default = "xoxb-xxx..."
}

variable "slack_signing_secret" {
  type = string
  default = ""
}