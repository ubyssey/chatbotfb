package main

import (
	"os"

	"github.com/ubyssey/chatbotfb/configuration"
)

func main() {
	mongoDBUrl := os.Getenv("MONGODBURL")
	if mongoDBUrl == "" {
		mongoDBUrl = configuration.Config.Database.MongoDB.LocalURL
	}

	// Connect to MongoDB
	database.Connect(mongoDBUrl)

	// Run HTTP Server
	server.Run()
}
