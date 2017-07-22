package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	campaignCollection  *mgo.Collection
	dbName              string
	userCollection      *mgo.Collection
	userCollectionError error
	senderID            string
	startCampaignId     = "6z479nb9-3x2f-23gs-g2dz-abc10625xc68"
)

// Creates a user record in MongoDB if non-existing user, otherwise
// update the user record
func createOrUpdateUser(opts messenger.MessageOpts) {
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

// Sends a reply back to the user depending on their message content
func handleReplyMessage(msg messenger.ReceivedMessage) {
	messageText := strings.ToLower(msg.Text)

	switch {
	case messageText == "menu":
		showMenu()
	case messageText == "start":
		startCampaign()
	default:
		chatbot.DefaultMessage(senderID)
	}
}

func initConfigVariables() {
	dbName = configuration.Config.Database.MongoDB.Name
	campaignCollection = database.MongoSession.DB(dbName).C("campaigns")
	userCollection = database.MongoSession.DB(dbName).C("users")
	// Check to see if the sender exists in the user collection
	userCollectionError = userCollection.FindId(senderID)
}

// Handles the POST message request from facebook
func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	senderID = opts.Sender.ID
	printlogger.Log("Received message from %s", senderID)

	// fetches the sender profile from facebook's Graph API
	_, profileErr := chatbot.CbMessenger.GetProfile(senderID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		printlogger.Log(profileErr.Error())
		return
	}

	initConfigVariables()
	createOrUpdateUser(opts)
	handleReplyMessage(msg)
}

func showMenu() {

}

// Start a new campaign for the user
func startCampaign() {
	startCampaign := campaign.Campaign{}
	startCampaignMissingErr := campaignCollection.FindId(startCampaignId).One(&startCampaign)

	if startCampaignMissingErr != nil {
		chatbot.DefaultMessage(senderID, "A start campaign was not found.")
		return
	}

	// Get the button templates to be shown to the user
	buttonsOptions, buttonOptionsErr := chatbot.GetButtonTemplateOptions(
		startCampaign.UUID,
		startCampaign.Nodes[startCampaign.RootNode].UserActions,
	)

	if buttonOptionsErr != nil {
		printlogger.Log(buttonOptionsErr.Error())
		return
	}

	// Initialize a message query
	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	mq.Template(
		template.GenericTemplate{
			Title:   startCampaign.Name,
			Buttons: buttonsOptions,
		},
	)

	resp, msgErr := chatbot.CbMessenger.SendMessage(mq)
	if msgErr != nil {
		printlogger.Log(msgErr.Error())
	}
	printlogger.Log("%+v", resp)
}
