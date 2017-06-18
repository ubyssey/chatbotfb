package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	chatbot.CbMessenger.MessageReceived = postMessage
}

// Handles the POST message request from facebook
func postMessage(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	// fetches the sender profile from facebook's Graph API

	_, profileErr := chatbot.CbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		fmt.Println(profileErr)
		return
	}

	// TODO: make the db stuff into a function. Ex. insertUser(db *mgo.Session ...). Also store user data?
	// TODO: make the bson fields more consistent
	// User collection (for MongoDB)
	dbName := configuration.Config.Database.MongoDB.Name

	uc := database.MongoSession.DB(dbName).C("users")
	user := models.User{}

	userCollectionError := uc.FindId(opts.Sender.ID).One(&user)

	// Check whether a user exists or not. If they are a first time user, create a record in database
	// otherwise update the record of that user
	if userCollectionError == nil {
		// existing user (user is found)

		set := bson.M{
			"lastSeen": time.Unix(opts.Timestamp, 0),
			"lastmessage": models.LastMessage{
				time.Now(),
				models.Event{
					"node",
					"4722d250-6162-4f02-a358-a4d55e3c8e20",
					"Nicasdfasdfasdfasde to meet you!",
				},
			},
		}

		uc.UpdateId(opts.Sender.ID, bson.M{"$set": set})

		printlogger.Log("Updated User %s", opts.Sender.ID)
	} else {
		// create new user

		uc.Insert(
			models.User{
				opts.Sender.ID,
				time.Unix(opts.Timestamp, 0),
				models.LastMessage{
					time.Now(),
					models.Event{
						"node",
						"4722d250-6162-4f02-a358-a4d55e3c8e20",
						"Nice to meet you!",
					},
				},
			},
		)

		printlogger.Log("Created User %s", opts.Sender.ID)
	}

	// Update the user activity timestamp

	if strings.ToLower(msg.Text) == "start" {
		mq := messenger.MessageQuery{}
		mq.RecipientID(opts.Sender.ID)
		mq.Template(template.GenericTemplate{
			Title: "abc",
			Buttons: []template.Button{
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: "test",
					Title:   "Nice to meet you!",
				},
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: "test",
					Title:   "I like NYT more than chatbots",
				},
			},
		})

		resp, err := chatbot.CbMessenger.SendMessage(mq)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v", resp)
	} else if len(msg.Text) > 0 {
		// chatbot only understands the message "start", any other message that is not a button or "start"
		// is invalid
		resp, err := chatbot.CbMessenger.SendSimpleMessage(
			opts.Sender.ID,
			fmt.Sprintf("Sorry, I don't understand your message."),
		)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v", resp)
	}

}
