package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type User struct {
	userID   string    `json:"userID" bson:"_id,omitempty"`
	lastSeen time.Time `json:"lastSeen" bson:"lastSeen"`
	LastMessage
}

type LastMessage struct {
	timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Event
}

type Event struct {
	nodeType string `json:"type" bson:"type"`
	target   string `json:"target" bson:"target"`
	label    string `json:"label" bson:"label"`
}

var (
	dbName = "chatbot"
)

func test() {

}
