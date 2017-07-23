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

// Sends a reply back to the user depending on their message content
func handleReplyMessage(msg messenger.ReceivedMessage) {
	messageText := strings.ToLower(msg.Text)

	switch {
	case messageText == "menu":
		messageoptions.ShowMenu()
	case messageText == "start":
		messageoptions.StartCampaign()
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
	user.createOrUpdateUser(opts)
	handleReplyMessage(msg)
}
