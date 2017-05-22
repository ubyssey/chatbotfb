package models

import (
	// "fmt"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "log"
	"time"
)

type User struct {
	UserID   string    `json:"userID" bson:"_id,omitempty"`
	LastSeen time.Time `json:"lastSeen" bson:"lastSeen"`
	LastMessage
}

type LastMessage struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Event
}

type Event struct {
	NodeType string `json:"type" bson:"type"`
	Target   string `json:"target" bson:"target"`
	Label    string `json:"label" bson:"label"`
}

func test() {

}
