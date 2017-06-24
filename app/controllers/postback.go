package controllers

import (
	"fmt"
	
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

func init() {
	chatbot.CbMessenger.Postback = postPostback
}

func postPostback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {
	_, profileErr := chatbot.CbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		printlogger.Log(profileErr)
		return

	}

	// Get the database name from the config
	dbName := configuration.Config.Database.MongoDB.Name
	// MongoDB campaign collection
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	userCollection := database.MongoSession.DB(dbName).C("users")

	postBackStruct := postback.Postback{}
	err := jsonparser.Parse(pb.Payload, postBackStruct)

	if err != nil {
		printlogger.Log("Error parsing the payload for user profile: %s", opts.Sender.ID)
		return
	}

	currentCampaign := campaignCollection.FindId(postBackStruct.CampaignId)

	if currentCampaign != nil {
		printLogger.Log("Error finding the campaign :%s", postBackStruct.CampaignId)
		return
	}

	if campaignNode, ok := currentCampaign[postBackStruct.CampaignId] {
		// If a node still has children, send a message with those children node options, 
		// otherwise send the final message of the current campaign
		if len(campaignNode.UserActions) > 0 {
			mq := messenger.MessageQuery{}
			mq.RecipientID(opts.Sender.ID)

			// Assume every node has two user actions
			firstPayloadOption := payload.Payload{
				CampaignId: postBackStruct.CampaignId,
				Event: &user.Event{
					NodeType: campaignNode.UserActions[0].NodeType,
					Target: campaignNode.UserActions[0].Target,
					Label: campaignNode.UserActions[0].Label,
				}
			}

			secondPayloadOption := payload.Payload{
				CampaignId: postBackStruct.CampaignId,
				Event: &user.Event{
					NodeType: campaignNode.UserActions[1].NodeType,
					Target: campaignNode.UserActions[1].Target,
					Label: campaignNode.UserActions[1].Label,
				}
			}

			// TODO: handle errors
			firstPayloadString, firstPayloadErr := jsonparser.ToJsonString(firstPayloadOption)
			secondPayloadString, secondPayloadErr := jsonparser.ToJsonString(secondPayloadOption)


			mq.Template(template.GenericTemplate{
				Title: c.Name,
				Buttons: []template.Button{
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: jsonparser.ToJsonString(firstPayloadOption),
						Title:   c.Nodes[c.RootNode].UserActions[0].Label,
					},
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: jsonparser.ToJsonString(secondPayloadOption),
						Title:   c.Nodes[c.RootNode].UserActions[1].Label,
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
		printlogger.Log("Campaign ID %s not found", postBackStruct.campaignId)
	}
}
