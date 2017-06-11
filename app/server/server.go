package server

import (
	"log"
	"net/http"
	"os"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"github.com/ubyssey/chatbotfb/app/controllers"
	"github.com/ubyssey/chatbotfb/app/routes"
	"github.com/ubyssey/chatbotfb/configuration"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	// Messenger SDK
	CbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
)

func Run() {
	CbMessenger.MessageReceived = controllers.GetMessage

	http.Handle("/", routes.Routes())

	// Heroku has its own env PORT. If not available use 3001 (for local development)
	port := os.Getenv("PORT")

	if port == "" {
		port = configuration.Database.MongoDB.LocalPort
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
