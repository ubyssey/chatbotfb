package controllers

import (
	"fmt"
	"net/http"

	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/configuration"
)

// Handles incoming requests for the /campaign endpoint
func GetCampaign(rw http.ResponseWriter, req *http.Request) {
	// TODO: query through every user in the database and send the campaign to everyone
	// TODO: have a table in the database to keep track of which campaign (by ID) were sent out
	// so that the same campaign doesn't get sent out multiple times. If an editor needs to edit
	// a campaign, the Mgmt-Api should create a new ID. Still need to think this through. What if
	// an editor needs to change a campaign while a user is in the middle of a campaign?

	var c campaign.Campaign
	if req.Body == nil {
		http.Error(rw, "Empty request body", 400)
		return
	}

	err := json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	// Get the database name from the config
	dbName := configuration.Config.Database.MongoDB.Name
	// MongoDB campaign collection
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	userCollection := database.MongoSession.DB(dbName).C("users")

	campaignCollectionError := campaignCollection.FindId(c.Name)

	// Check whether or not the campaign exists or not. If it does not exist, add it to the
	// database. If it does, then do nothing.
	if campaignCollectionError != nil {
		// Campaign does not exist

		campaignCollection.Insert(c)

		// TODO: for every user, update their last message once a new campaign is sent over.

	}
}
