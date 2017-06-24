package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	"gopkg.in/mgo.v2/bson"
)

// Sends a GET request to the mgmt API
// TOOD: implement an actual HTTP request once the mgmt API endpoint is implemented
func GetCampaignFromMgmtApi() {
	testCampaignNodePath, _ := filepath.Abs("campaign-node.json")
	raw, err := ioutil.ReadFile(testCampaignNodePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var newCampaign campaign.Campaign
	json.Unmarshal(raw, &newCampaign)

	// Get the database name from the config
	dbName := configuration.Config.Database.MongoDB.Name
	// MongoDB campaign collection
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	userCollection := database.MongoSession.DB(dbName).C("users")

	campaignCollectionError := campaignCollection.FindId(newCampaign.UUID)

	// Check whether or not the campaign exists or not. If it does not exist, add it to the
	// database. If it does, then do nothing.
	if campaignCollectionError != nil {
		// Campaign does not exist so insert a new campaign
		campaignCollection.Insert(newCampaign)

		u := user.User{}

		// for every user, update their last message once a new campaign is sent over.
		// TODO: find a way to use updateAll instead. Also refactor it later.
		findUsers := userCollection.Find(bson.M{})
		users := findUsers.Iter()
		for users.Next(&u) {
			set := bson.M{
				"lastmessage": &user.LastMessage{
					time.Now(),
					user.Event{
						newCampaign.Nodes[newCampaign.RootNode].UserActions[0].NodeType,
						newCampaign.Nodes[newCampaign.RootNode].UserActions[0].Target,
						newCampaign.Nodes[newCampaign.RootNode].UserActions[0].Label,
					},
				},
			}
			userCollection.UpdateId(u.UserID, bson.M{"$set": set})

			_, profileErr := chatbot.CbMessenger.GetProfile(u.UserID)
			// if the sender profile is invalid, print out error and return
			if profileErr != nil {
				printlogger.Log(profileErr.Error())
				return

			}

			// Assume for a new campaign, there is atleast a root node with two user actions

			mq := messenger.MessageQuery{}
			mq.RecipientID(u.UserID)

			firstPayloadOption := payload.Payload{
				CampaignId: newCampaign.UUID,
				Event: &user.Event{
					NodeType: newCampaign.Nodes[newCampaign.RootNode].UserActions[0].NodeType,
					Target:   newCampaign.Nodes[newCampaign.RootNode].UserActions[0].Target,
					Label:    newCampaign.Nodes[newCampaign.RootNode].UserActions[0].Label,
				},
			}

			secondPayloadOption := payload.Payload{
				CampaignId: newCampaign.UUID,
				Event: &user.Event{
					NodeType: newCampaign.Nodes[newCampaign.RootNode].UserActions[1].NodeType,
					Target:   newCampaign.Nodes[newCampaign.RootNode].UserActions[1].Target,
					Label:    newCampaign.Nodes[newCampaign.RootNode].UserActions[1].Label,
				},
			}

			// TODO: handle errors
			firstPayloadString, firstPayloadErr := jsonparser.ToJsonString(firstPayloadOption)
			secondPayloadString, secondPayloadErr := jsonparser.ToJsonString(secondPayloadOption)

			// TODO: handle errors
			if firstPayloadErr != nil {

			}

			if secondPayloadErr != nil {

			}

			mq.Template(template.GenericTemplate{
				Title: newCampaign.Name,
				Buttons: []template.Button{
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: firstPayloadString,
						Title:   newCampaign.Nodes[newCampaign.RootNode].UserActions[0].Label,
					},
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: secondPayloadString,
						Title:   newCampaign.Nodes[newCampaign.RootNode].UserActions[1].Label,
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
}
