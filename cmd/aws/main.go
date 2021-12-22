package main

import (
	handler "github.com/akubi0w1/chatbot-202201/external/handler/aws"
	"github.com/akubi0w1/chatbot-202201/log"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	logger := log.New()
	app := handler.NewApplication()

	logger.Infof("lambda start")
	lambda.Start(app.Handle)
}
