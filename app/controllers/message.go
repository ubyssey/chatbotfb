package controllers

import (
	"fmt"
	"strings"
	"time"

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
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	dbName   string
	uc       *mgo.Collection // User Collection
	ucError  error           // User Collection error
	senderID string
)

// Creates a user record in MongoDB if non-existing user, otherwise
// update the user record
func createOrUpdateUser(opts messenger.MessageOpts) {
	// Check whether a user exists or not. If they are a first time user, create a record in database
	// otherwise update the record of that user
	// TODO: figure out when lastMessage is updated since there are two user actions
	if ucError == nil {
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

		uc.UpdateId(senderID, bson.M{"$set": set})

		printlogger.Log("Updated User %s", senderID)
	} else {
		// create new user
		uc.Insert(
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
func handleReplyMessage(opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")

	if strings.ToLower(msg.Text) == "start" {
		// A new starting campaign for the user
		startCampaign := campaign.Campaign{}
		startCampaignMissingErr := campaignCollection.FindId("6z479nb9-3x2f-23gs-g2dz-abc10625xc68").One(&startCampaign)

		// TODO: handle err
		if startCampaignMissingErr != nil {
			startCampaignErrResponse, _ := chatbot.CbMessenger.SendSimpleMessage(
				opts.Sender.ID,
				fmt.Sprintf("A start campaign was not found."),
			)

			printlogger.Log("%+v", startCampaignErrResponse)
			return
		}

		// Assume every starting campaign node has two user actions
		firstPayloadOption := payload.Payload{
			CampaignId: startCampaign.UUID,
			Event: &user.Event{
				NodeType: startCampaign.Nodes[startCampaign.RootNode].UserActions[0].NodeType,
				Target:   startCampaign.Nodes[startCampaign.RootNode].UserActions[0].Target,
				Label:    startCampaign.Nodes[startCampaign.RootNode].UserActions[0].Label,
			},
		}

		secondPayloadOption := payload.Payload{
			CampaignId: startCampaign.UUID,
			Event: &user.Event{
				NodeType: startCampaign.Nodes[startCampaign.RootNode].UserActions[1].NodeType,
				Target:   startCampaign.Nodes[startCampaign.RootNode].UserActions[1].Target,
				Label:    startCampaign.Nodes[startCampaign.RootNode].UserActions[1].Label,
			},
		}

		firstPayloadString, firstPayloadErr := jsonparser.ToJsonString(firstPayloadOption)
		secondPayloadString, secondPayloadErr := jsonparser.ToJsonString(secondPayloadOption)

		// TODO: handle errors
		if firstPayloadErr != nil {

		}

		if secondPayloadErr != nil {

		}

		// Initialize a message query
		mq := messenger.MessageQuery{}
		mq.RecipientID(senderID)

		mq.Template(template.GenericTemplate{
			Title: startCampaign.Name,
			Buttons: []template.Button{
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: firstPayloadString,
					Title:   startCampaign.Nodes[startCampaign.RootNode].UserActions[0].Label,
				},
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: secondPayloadString,
					Title:   startCampaign.Nodes[startCampaign.RootNode].UserActions[1].Label,
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

// Handles the POST message request from facebook
func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	printlogger.Log("Received message from %s", opts.Sender.ID)

	// fetches the sender profile from facebook's Graph API
	_, profileErr := chatbot.CbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		fmt.Println(profileErr)
		return
	}

	senderID = opts.Sender.ID
	dbName = configuration.Config.Database.MongoDB.Name

	// User collection (for Mon
	uc = database.MongoSession.DB(dbName).C("users")
	currUser := user.User{}
	ucError = uc.FindId(senderID).One(&currUser)

	createOrUpdateUser(opts)
	handleReplyMessage(opts, msg)
}
