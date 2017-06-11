package database

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoSession    *mgo.Session
	mongoSessionErr error
)

func Connect(databaseURL string) {
	mongoSession, mongoSessionErr = mgo.Dial(databaseURL)

	// panic if mongoDB session fails initialization
	if mongoSessionErr != nil {
		panic(mongoSessionErr)
	}
	defer mongoSession.Close()
}
