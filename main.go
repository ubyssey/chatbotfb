package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var cbMessenger = &messenger.Messenger{
	AccessToken: os.Getenv("TOKEN"),
	Debug:       cbMessenger.DebugAll,
}

func main() {
	cbMessenger.MessageReceived = MessageReceived
	http.HandlerFunc("/webhook", cbMessenger.Handler)
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := cbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if err != nil {
		fmt.Println(err)
		return
	}

	// send a simple message
	resp, err := mess.SendSimpleMessage(
		opts.Sender.ID, 
		fmt.Sprintf("Hello, %s %s, %s", profile.FirstName, profile.LastName, msg.Text)
	)

	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Printf("%+v", resp)
}