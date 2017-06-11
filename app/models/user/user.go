package user

import (
	"time"
)

// TODO: make the BSON fields consistent
type User struct {
	UserID      string    `json:"userID" bson:"_id,omitempty"`
	LastSeen    time.Time `json:"lastSeen" bson:"lastSeen"`
	LastMessage `json:"lastMessage" bson:"lastMessage"`
}

type LastMessage struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Event     `json:"event" bson:"event"`
}

type Event struct {
	NodeType string `json:"type" bson:"type"`
	Target   string `json:"target" bson:"target"`
	Label    string `json:"label" bson:"label"`
}
