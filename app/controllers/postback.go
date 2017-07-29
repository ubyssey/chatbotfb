package controllers

import (
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/messageactions"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

// Sends back a reply when a user clicks a 'postback' button
func Postback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {
	senderID := opts.Sender.ID
	_, profileErr := chatbot.CbMessenger.GetProfile(senderID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		printlogger.Log(profileErr.Error())
		return
	}

	// Get the payload from the postback message
	payloadStruct, payloadStructErr := payload.GetPayloadStruct(pb)
	if payloadStructErr != nil {
		printlogger.Log(payloadStructErr.Error())
		printlogger.Log("Error parsing the payload '%s' for user profile: %s", pb.Payload, senderID)
		return
	}

	switch {
	case payloadStruct.UrlLink != "":
		messageactions.ShowMenuListTemplate(senderID, payloadStruct)
	case payloadStruct.CampaignId != "":
		messageactions.SendNextCampaignNodeActions(senderID, payloadStruct)
	default:
	}
}
