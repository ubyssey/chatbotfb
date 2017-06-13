package server

import (
	"log"
	"net/http"
	"os"

	"github.com/ubyssey/chatbotfb/app/routes"
	"github.com/ubyssey/chatbotfb/configuration"
)

func Run() {
	http.Handle("/", routes.Routes())

	// Heroku has its own env PORT. If not available use 3001 (for local development)
	port := os.Getenv("PORT")

	if port == "" {
		port = configuration.Config.Database.MongoDB.LocalPort
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
