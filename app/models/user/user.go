package user

import (
	"time"
)

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

// Creates a user record in MongoDB if non-existing user, otherwise
// update the user record
func CreateOrUpdateUser(opts messenger.MessageOpts) {
	// Check whether a user exists or not. If they are a first time user, create a record in database
	// otherwise update the record of that user
	// TODO: figure out when lastMessage is updated since there are two user actions
	if userCollectionError == nil {
		// existing user (user is found)
		set := bson.M{
			"lastSeen": time.Unix(opts.Timestamp, 0),
			"lastMessage": &user.LastMessage{
				time.Now(),
				user.Event{
					"node",
					"4722d250-6162-4f02-a358-a4d55e3c8e20",
					"Nicasdfasdfasdfasde to meet you!",
				},
			},
		}

		userCollection.UpdateId(senderID, bson.M{"$set": set})

		printlogger.Log("Updated User %s", senderID)
	} else {
		// create new user
		userCollection.Insert(
			&user.User{
				senderID,
				time.Unix(opts.Timestamp, 0),
				user.LastMessage{
					time.Now(),
					user.Event{
						"node",
						"4722d250-6162-4f02-a358-a4d55e3c8e20",
						"Nice to meet you!",
					},
				},
			},
		)

		printlogger.Log("Created User %s", senderID)
	}
}
