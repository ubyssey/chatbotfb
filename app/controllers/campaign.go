package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Sends a GET request to the mgmt API
// TOOD: implement an actual HTTP request once the mgmt API endpoint is implemented
func GetCampaignFromMgmtApi() {
	raw, err := ioutil.ReadFile("../../../campaign-node.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c campaign.Campaign
	json.Unmarshal(raw, &campaign)

	// Get the database name from the config
	dbName := configuration.Config.Database.MongoDB.Name
	// MongoDB campaign collection
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	userCollection := database.MongoSession.DB(dbName).C("users")

	campaignCollectionError := campaignCollection.FindId(c.UUID)

	// Check whether or not the campaign exists or not. If it does not exist, add it to the
	// database. If it does, then do nothing.
	if campaignCollectionError != nil {
		// Campaign does not exist so insert a new campaign
		campaignCollection.Insert(c)

		u := user.User

		// for every user, update their last message once a new campaign is sent over.
		// TODO: find a way to use updateAll instead. Also refactor it later.
		findUsers := userCollection.Find(bson.M{})
		users := findUsers.Next(&u) {
			set := bson.M{
				"lastmessage": user.LastMessage{
					time.Now(),
					user.Event{
						c.Nodes[c.RootNode].UserActions[0].NodeType
						c.Nodes[c.RootNode].UserActions[0].Target,
						c.Nodes[c.RootNode].UserActions[0].Label,
					},
				},
			}
			userCollection.UpdateId(u.USERID, bson.M{"$set": set})

			// send a campaign message to the user
			// TODO: payload can only be a string. Make it a JSON string instead?
			mq := messenger.MessageQuery{}
			mq.RecipientID(u.USERID)
			mq.Template(template.GenericTemplate{
				Title: c.Name,
				Buttons: []template.Button{
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: c.Nodes[c.RootNode].UserActions[0].Target,
						Title:   c.Nodes[c.RootNode].UserActions[0].Label,
					},
					template.Button{
						Type:    template.ButtonTypePostback,
						Payload: c.Nodes[c.RootNode].UserActions[1].Target,
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
}
