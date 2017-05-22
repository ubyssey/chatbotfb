package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type User struct {
	userID   string
	lastSeen time.Time
	LastMessage
}

type LastMessage struct {
	timestamp time.Time
	Event
}

type Event struct {
	nodeType string
	target   string
	label    string
}

var (
	dbName = "chatbot"
)

func test() {

}
