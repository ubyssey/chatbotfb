package controllers

import (
	"fmt"
	"net/http"
)

// Handles incoming requests for the /campaign endpoint
func GetCampaign(rw http.ResponseWriter, req *http.Request) {
	// TODO: query through every user in the database and send the campaign to everyone
	// TODO: have a table in the database to keep track of which campaign (by ID) were sent out
	// so that the same campaign doesn't get sent out multiple times. If an editor needs to edit
	// a campaign, the Mgmt-Api should create a new ID. Still need to think this through. What if
	// an editor needs to change a campaign while a user is in the middle of a campaign?

	fmt.Fprintf(rw, "<h1>Campaign</h1>")
}
