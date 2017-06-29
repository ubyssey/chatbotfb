package server

import (
	"log"
	"net/http"
	"os"

	"github.com/ubyssey/chatbotfb/app/controllers"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/routes"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"
)

func Run() {
	chatbot.CbMessenger.MessageReceived = controllers.MessageReceived
	chatbot.CbMessenger.Postback = controllers.Postback

	http.Handle("/", routes.Routes())

	// Testing a "would-be" push notification from mgmt-api
	controllers.GetCampaignFromMgmtApi()

	// Heroku has its own env PORT. If not available use 3001 (for local development)
	port := os.Getenv("PORT")

	if port == "" {
		port = configuration.Config.Database.MongoDB.LocalPort
	}

	printlogger.Log("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
