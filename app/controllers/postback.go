package controllers

import (
	"fmt"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	senderID string
)

// Sends back a reply when a user clicks a 'postback' button
func Postback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {
	senderID = opts.Sender.ID
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

	// Get the campaign from the postback's 'campaignId' field
	dbName := configuration.Config.Database.MongoDB.Name
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	currentCampaign, currentCampaignErr := campaign.GetCampaignStruct(campaignCollection, payloadStruct.CampaignId)
	if currentCampaignErr != nil {
		printlogger.Log(currentCampaignErr.Error())
		printlogger.Log("Error finding the campaign :%s", payloadStruct.CampaignId)
		return
	}

	sendNextNodeActions(currentCampaign, payloadStruct)
}

// Send the user actions for the next node
func sendNextNodeActions(currentCampaign campaign.Campaign, payloadStruct payload.Payload) {
	if campaignNode, ok := currentCampaign.Nodes[payloadStruct.Event.Target]; ok {
		// If a node still has children, send a message with those children node options,
		// otherwise send the final message of the current campaign
		if len(campaignNode.UserActions) > 0 {
			mq := messenger.MessageQuery{}
			mq.RecipientID(senderID)

			buttonsOptions, buttonOptionsErr := chatbot.GetButtonTemplateOptions(
				payloadStruct.CampaignId,
				campaignNode.UserActions,
			)

			if buttonOptionsErr != nil {
				printlogger.Log(buttonOptionsErr.Error())
				return
			}

			// Generic Message Query template to be sent to the user
			mq.Template(
				template.GenericTemplate{
					Title:   currentCampaign.Name,
					Buttons: buttonsOptions,
				},
			)

			resp, msgErr := chatbot.CbMessenger.SendMessage(mq)

			if msgErr != nil {
				printlogger.Log(msgErr.Error())
			}

			printlogger.Log("%+v", resp)
		} else {
			chatbot.DefaultMessage(senderID)
		}
	} else {
		printlogger.Log("Campaign Node target %s not found for user %s", payloadStruct.Event.Target, senderID)
	}
}
