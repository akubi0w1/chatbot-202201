package main

import (
	"fmt"
	"net/http"
	"os"

	handler "github.com/akubi0w1/chatbot-202201/external/handler/gcp"
	"github.com/akubi0w1/chatbot-202201/log"
)

func main() {
	logger := log.New()
	app := handler.NewApplication()
	mux := app.Routing()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Infof("start server port=%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		logger.Fatalf("failed to http listen and serve: %v", err)
	}
}
