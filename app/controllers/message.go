package controllers

import (
	"strings"

	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/messageactions"
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Sends a reply back to the user depending on their message content
func handleReplyMessage(opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	messageText := strings.ToLower(msg.Text)

	switch {
	case messageText == "menu":
		messageactions.ShowMenu()
	case messageText == "start":
		messageactions.StartCampaign(opts.Sender.ID)
	default:
		chatbot.DefaultMessage(opts.Sender.ID, "")
	}
}

// Handles the POST message request from facebook
func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	// fetches the sender profile from facebook's Graph API
	_, profileErr := chatbot.CbMessenger.GetProfile(senderID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		printlogger.Log(profileErr.Error())
		return
	}

	printlogger.Log("Received message from %s", opts.Sender.ID)

	user.createOrUpdateUser(opts)
	handleReplyMessage(opts, msg)
}
