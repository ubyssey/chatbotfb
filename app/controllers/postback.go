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
		printlogger.Log("Error parsing the payload")
		return
	}

	campaignCollection := campaignCollection.FindId(postBackStruct.campaignId)

	if campaignCollection != nil {
		printLogger.Log("Error finding the campaign :%s", postBackStruct.campaignId)
		return
	}

	if campaignNode, ok := campaignCollection[postBackStruct.campaignId] {
		mq := messenger.MessageQuery{}
		mq.RecipientID(opts.Sender.ID)

		// TODO: implement the payload options. Not every node has a payload / user action
		firstPayloadOption := payload.Payload{

		}

		secondPayloadOption := payload.Payload{

		}

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
	}
}
