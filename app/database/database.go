package database

import (
	"gopkg.in/mgo.v2"

	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
)

var (
	MongoSession    *mgo.Session
	mongoSessionErr error
)

func Connect(databaseURL string) {
	MongoSession, mongoSessionErr = mgo.Dial(databaseURL)

	printlogger.Log("Connecting to MongoDB at : %s", databaseURL)

	// panic if mongoDB session fails initialization
	if mongoSessionErr != nil {
		panic(mongoSessionErr)
	}
}

func Disconnect() {
	MongoSession.Close()
}
