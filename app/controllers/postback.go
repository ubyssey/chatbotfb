package controllers

import (
	"fmt"

	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

func Postback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {
	_, profileErr := chatbot.CbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		printlogger.Log(profileErr.Error())
		return

	}

	// Get the database name from the config
	dbName := configuration.Config.Database.MongoDB.Name
	// MongoDB campaign collection
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")

	payloadStruct := payload.Payload{}
	err := jsonparser.Parse([]byte(pb.Payload), payloadStruct)

	if err != nil {
		printlogger.Log("Error parsing the payload for user profile: %s", opts.Sender.ID)
		return
	}

	currentCampaign := campaign.Campaign{}
	currentCampaignQuery := campaignCollection.FindId(payloadStruct.CampaignId).One(&currentCampaign)

	if currentCampaignQuery != nil {
		printlogger.Log("Error finding the campaign :%s", payloadStruct.CampaignId)
		return
	}

	if campaignNode, ok := currentCampaign.Nodes[payloadStruct.CampaignId]; ok {
		// If a node still has children, send a message with those children node options,
		// otherwise send the final message of the current campaign
		if len(campaignNode.UserActions) > 0 {
			mq := messenger.MessageQuery{}
			mq.RecipientID(opts.Sender.ID)

			// Assume every node has two user actions
			firstPayloadOption := payload.Payload{
				CampaignId: payloadStruct.CampaignId,
				Event: &user.Event{
					NodeType: campaignNode.UserActions[0].NodeType,
					Target:   campaignNode.UserActions[0].Target,
					Label:    campaignNode.UserActions[0].Label,
				},
			}

			secondPayloadOption := payload.Payload{
				CampaignId: payloadStruct.CampaignId,
				Event: &user.Event{
					NodeType: campaignNode.UserActions[1].NodeType,
					Target:   campaignNode.UserActions[1].Target,
					Label:    campaignNode.UserActions[1].Label,
				},
			}

			firstPayloadString, firstPayloadErr := jsonparser.ToJsonString(firstPayloadOption)
			secondPayloadString, secondPayloadErr := jsonparser.ToJsonString(secondPayloadOption)

			// TODO: handle errors
			if firstPayloadErr != nil {
				printlogger.Log("%s", firstPayloadErr.Error())
				return
			}

			if secondPayloadErr != nil {
				printlogger.Log("%s", secondPayloadErr.Error())
				return
			}

			mq.Template(template.GenericTemplate{
				Title: currentCampaign.Name,
				Buttons: []template.Button{
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: firstPayloadString,
						Title:   currentCampaign.Nodes[currentCampaign.RootNode].UserActions[0].Label,
					},
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: secondPayloadString,
						Title:   currentCampaign.Nodes[currentCampaign.RootNode].UserActions[1].Label,
					},
				},
			})

			resp, err := chatbot.CbMessenger.SendMessage(mq)

			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%+v", resp)
		} else {
			resp, err := chatbot.CbMessenger.SendSimpleMessage(
				opts.Sender.ID,
				fmt.Sprintf(campaignNode.Content.Text),
			)

			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%+v", resp)
		}
	} else {
		printlogger.Log("Campaign ID %s not found", payloadStruct.CampaignId)
	}
}
