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

	printlogger.Log("Received payload '%s' from user %s", pb.Payload, opts.Sender.ID)

	payloadStruct := payload.Payload{}
	err := jsonparser.Parse([]byte(pb.Payload), &payloadStruct)

	if err != nil {
		printlogger.Log(err.Error())
		printlogger.Log("Error parsing the payload '%s' for user profile: %s", pb.Payload, opts.Sender.ID)
		return
	}

	currentCampaign := campaign.Campaign{}
	currentCampaignQuery := campaignCollection.FindId(payloadStruct.CampaignId).One(&currentCampaign)

	if currentCampaignQuery != nil {
		printlogger.Log("Error finding the campaign :%s", payloadStruct.CampaignId)
		return
	}

	if campaignNode, ok := currentCampaign.Nodes[payloadStruct.Event.Target]; ok {
		// If a node still has children, send a message with those children node options,
		// otherwise send the final message of the current campaign
		if len(campaignNode.UserActions) > 0 {
			mq := messenger.MessageQuery{}
			mq.RecipientID(opts.Sender.ID)

			// A button slice to hold each button option to be shown to the user
			buttonsSlice := []template.Button{}
			var button template.Button

			for _, currUserAction := range campaignNode.UserActions {
				// Reset the button struct every loop
				button = template.Button{}

				// If the node type is a "link", create a NewWebURLButton template. Otherwise,
				// if it is a "node", then create a postback Button template with its payload
				if currUserAction.NodeType == "link" {
					button = template.NewWebURLButton(
						currUserAction.Label,
						currUserAction.Target
					)
				} else if currUserAction.NodeType == "node" {
					payloadOption := payload.Payload{
						CampaignId: payloadStruct.CampaignId,
						Event: &user.Event{
							NodeType: currUserAction.NodeType,
							Target:   currUserAction.Target,
							Label:    currUserAction.Label,
						},
					}

					payloadOptionString, payloadOptionParsingErr := jsonparser.ToJsonString(payloadOption)

					if payloadOptionParsingErr != nil {
						printlogger.Log("%s", payloadOptionParsingErr.Error())
						return
					}

					button = template.Button{
						Type:    template.ButtonTypePostback,
						Payload: payloadOptionString,
						Title:   currUserAction.Label,
					}
				}

				buttonsSlice = append(buttonsSlice, button)
			}

			// Generic Message Query template to be sent to the user
			mq.Template(
				template.GenericTemplate{
					Title: currentCampaign.Name,
					Buttons: buttonsSlice,
				}
			)

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
		printlogger.Log("Campaign Node target %s not found for user %s", payloadStruct.Event.Target, opts.Sender.ID)
	}
}
