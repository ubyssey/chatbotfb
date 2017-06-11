package server

import (
	"log"
	"net/http"
	"os"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	// Messenger SDK
	cbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
)

func Run() {
	cbMessenger.MessageReceived = MessageReceived

	// TODO: handle routes seperately in another file
	// API endpoints
	http.HandleFunc("/webhook", cbMessenger.Handler)
	http.HandleFunc("/campaign", campaignHandler)

	// Heroku has its own env PORT. If not available use 3001 (for local development)
	port := os.Getenv("PORT")
	if port == "" {
		port = config.Database.MongoDB.LocalPort
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
