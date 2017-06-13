package database

import (
	"gopkg.in/mgo.v2"
)

var (
	MongoSession    *mgo.Session
	mongoSessionErr error
)

func Connect(databaseURL string) {
	MongoSession, mongoSessionErr = mgo.Dial(databaseURL)

	// panic if mongoDB session fails initialization
	if mongoSessionErr != nil {
		panic(mongoSessionErr)
	}
	defer MongoSession.Close()
}
