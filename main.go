package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

type campaignNode struct {}

var cbMessenger = &messenger.Messenger{
	AccessToken: os.Getenv("TOKEN"),
}

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := cbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if err != nil {
		fmt.Println(err)
		return
	}

	// send a simple message
	resp, err := cbMessenger.SendSimpleMessage(
		opts.Sender.ID,
		fmt.Sprintf("Hello, %s %s, %s", profile.FirstName, profile.LastName, msg.Text),
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", resp)
}

// Handles incoming requests for the /campaign endpoint
func campaignHandler(rw http.ResponseWriter, req *http.Request) {
	file, err := ioutil.ReadFile("./campaign-node.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// TODO: query through every user in the database and send the campaign to everyone
	// TODO: have a table in the database to keep track of which campaign (by ID) were sent out
	// so that the same campaign doesn't get sent out multiple times. If an editor needs to edit
	// a campaign, the Mgmt-Api should create a new ID. Still need to think this through. What if
	// an editor needs to change a campaign while a user is in the middle of a campaign? 
	// TODO: ask at next meeting if its necessary to hash the details of users?


	fmt.Fprintf(rw, "<h1>Testing 123</h1><div>hey/div>")
}

func main() {
	// TODO: remove the TOKEN variable in production. It doesn't matter if the token is seen
	// since it is only a test cb page. This is also only used locally on your computer, 
	// Heroku supports environment variables. 
	apiToken := "EAATZAxfQTVYQBAE08rAVR3NZBaPwe0FTDEbGZBLbIKx5LUf5Y5m2DiZAkg1ZBhxo0IKhGuLHMkMj3ZAXdOJygZBTK9KZCyGb8J87uxsxpFXZCrzZByLveD2cmHryuCxDNtv2ifVlM18J1QoktcHLDaJI59Vlvn120t613QrQ2Ae0GSnwZDZD"
	os.SetEnv("TOKEN", apiToken)

	cbMessenger.MessageReceived = MessageReceived

	http.HandleFunc("/webhook", cbMessenger.Handler)
	http.HandleFunc("/campaign", campaignHandler)

	log.Fatal(http.ListenAndServe(":3001", nil))
}