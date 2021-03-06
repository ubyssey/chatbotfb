package user

import (
	"time"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2/bson"
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
	senderID := opts.Sender.ID
	dbName := configuration.Config.Database.MongoDB.Name
	userCollection := database.MongoSession.DB(dbName).C("users")
	// Check to see if the sender exists in the user collection
	userCollectionError := userCollection.FindId(senderID)

	// If the user is a first time user, create a record in database.
	// Otherwise update the record of that user.
	// TODO: figure out when lastMessage is updated since there are two user actions
	if userCollectionError == nil {
		// existing user (user is found)
		set := bson.M{
			"lastSeen": time.Unix(opts.Timestamp, 0),
			"lastMessage": &LastMessage{
				time.Now(),
				Event{
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
			&User{
				senderID,
				time.Unix(opts.Timestamp, 0),
				LastMessage{
					time.Now(),
					Event{
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
