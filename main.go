package main

import (
	"os"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/server"
	"github.com/ubyssey/chatbotfb/configuration"
)

func main() {
	mongoDBUrl := os.Getenv("MONGODBURL")
	if mongoDBUrl == "" {
		mongoDBUrl = configuration.Config.Database.MongoDB.LocalURL
	}

	// Connect to MongoDB
	database.Connect(mongoDBUrl)
	defer database.Disconnect()

	// Run HTTP Server
	server.Run()
}
