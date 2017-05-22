package main

import (
	// Standard packages / libraries
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	// Internal packages
	"models"

	// External packages
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// TODO: Combine these vars into a group?
var cbMessenger = &messenger.Messenger{
	AccessToken: os.Getenv("TOKEN"),
}

var dbName = "chatbot"

// Initialize mongoDB session
var mongoSession, mongoSessionErr = mgo.Dial("<INSERT MONGO DB INSTANCE URL HERE>")

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	// fetches the sender profile from facebook's Graph API
	profile, err := cbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: make the db stuff into a function. Ex. insertUser(db *mgo.Session ...)
	// User collection (for MongoDB)
	uc := session.DB(dbName).C("users")
	user := User{}

	err := uc.Find(bson.M{"userID": opts.Sender.ID}).One(&user)

	if err != nil {
		// existing user (user is found)
	} else {
		// new user

	}

	// Update the user activity timestamp

	if toLower(msg) == "start" {
		startTestCampaign()
		return
	} else if len(msg) > 0 {
		// chatbot only understands the message "start", any other message that is not a button or "start"
		// is invalid
		resp, err := cbMessenger.SendSimpleMessage(
			opts.Sender.ID,
			fmt.Sprintf("Sorry, I don't understand your message."),
		)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v", resp)
	}

}

// Handles incoming requests for the /campaign endpoint
func campaignHandler(rw http.ResponseWriter, req *http.Request) {
	// TODO: query through every user in the database and send the campaign to everyone
	// TODO: have a table in the database to keep track of which campaign (by ID) were sent out
	// so that the same campaign doesn't get sent out multiple times. If an editor needs to edit
	// a campaign, the Mgmt-Api should create a new ID. Still need to think this through. What if
	// an editor needs to change a campaign while a user is in the middle of a campaign?
	// TODO: ask at next meeting if its necessary to hash the details of users?

	fmt.Fprintf(rw, "<h1>Campaign</h1>")
}

func genericErrorLogger() {
	// TODO: find a method that logs the timestamp and the error or implement it yourself
	return
}

// Test campaign
func startTestCampaign() {
	raw, err := ioutil.ReadFile("./campaign-node.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	// panic if mongoDB session fails initialization
	if mongoSessionErr != nil {
		panic(mongoSessionErr)
	}
	defer mongoSession.close()

	// TODO: remove the TOKEN variable in production. It doesn't matter if the token is seen
	// since it is only a test cb page. This is also only used locally on your computer,
	// Heroku supports environment variables.
	apiToken := "EAATZAxfQTVYQBAE08rAVR3NZBaPwe0FTDEbGZBLbIKx5LUf5Y5m2DiZAkg1ZBhxo0IKhGuLHMkMj3ZAXdOJygZBTK9KZCyGb8J87uxsxpFXZCrzZByLveD2cmHryuCxDNtv2ifVlM18J1QoktcHLDaJI59Vlvn120t613QrQ2Ae0GSnwZDZD"
	os.Setenv("TOKEN", apiToken)

	cbMessenger.MessageReceived = MessageReceived

	// API endpoints
	http.HandleFunc("/webhook", cbMessenger.Handler)
	http.HandleFunc("/campaign", campaignHandler)

	// TODO: Set it to Heroku's env PORT || 3001 when deplying to Heroku
	log.Fatal(http.ListenAndServe(":3001", nil))
}
