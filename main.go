package main

import (
	// Standard packages / libraries
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// Internal packages
	"github.com/ubyssey/chatbotfb/models"

	// External packages
	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mongoSession    *mgo.Session
	mongoSessionErr error

	// Messenger SDK
	cbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
	dbName = "ubysseycb"
)

func MessageReceived(event messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	// fetches the sender profile from facebook's Graph API

	_, profileErr := cbMessenger.GetProfile(opts.Sender.ID)
	// if the sender profile is invalid, print out error and return
	if profileErr != nil {
		fmt.Println(profileErr)
		return
	}

	// TODO: make the db stuff into a function. Ex. insertUser(db *mgo.Session ...). Also store user data?
	// TODO: make the bson fields more consistent
	// User collection (for MongoDB)
	uc := mongoSession.DB(dbName).C("users")
	user := models.User{}

	userCollectionError := uc.FindId(opts.Sender.ID).One(&user)

	if userCollectionError == nil {
		// existing user (user is found)
		set := bson.M{
			"lastSeen": time.Unix(opts.Timestamp, 0),
			"lastmessage": &models.LastMessage{
				time.Now(),
				models.Event{
					"node",
					"4722d250-6162-4f02-a358-a4d55e3c8e20",
					"Nicasdfasdfasdfasde to meet you!",
				},
			},
		}

		uc.UpdateId(opts.Sender.ID, bson.M{"$set": set})

		PrintLogger("Updated User %s", opts.Sender.ID)
	} else {
		// create new user
		uc.Insert(
			&models.User{
				opts.Sender.ID,
				time.Unix(opts.Timestamp, 0),
				models.LastMessage{
					time.Now(),
					models.Event{
						"node",
						"4722d250-6162-4f02-a358-a4d55e3c8e20",
						"Nice to meet you!",
					},
				},
			},
		)

		PrintLogger("Created User %s", opts.Sender.ID)
	}

	// Update the user activity timestamp

	if strings.ToLower(msg.Text) == "start" {
		mq := messenger.MessageQuery{}
		mq.RecipientID(opts.Sender.ID)
		mq.Template(template.GenericTemplate{
			Title: "abc",
			Buttons: []template.Button{
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: "test",
					Title:   "Nice to meet you!",
				},
				template.Button{
					Type:    template.ButtonTypePostback,
					Payload: "test",
					Title:   "I like NYT more than chatbots",
				},
			},
		})

		resp, err := cbMessenger.SendMessage(mq)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v", resp)
	} else if len(msg.Text) > 0 {
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

// Generic logging with time prefix
func PrintLogger(format string, args ...interface{}) {
	// Include stack traces maybe? See errgo package
	fmt.Printf("[LOG] "+time.Now().Format("2017-05-27 00:00:00")+" "+format, args...)
}

func main() {
	// Initialize MongoDB session
	mongoDBUrl := os.Getenv("MONGODBURL")
	if mongoDBUrl == "" {
		mongoDBUrl = "mongodb://127.0.0.1:27017"
	}

	mongoSession, mongoSessionErr = mgo.Dial(mongoDBUrl)

	// panic if mongoDB session fails initialization
	if mongoSessionErr != nil {
		panic(mongoSessionErr)
	}
	defer mongoSession.Close()

	// TODO: remove the TOKEN variable in production. It doesn't matter if the token is seen
	// since it is only a test cb page. This is also only used locally on your computer,
	// Heroku supports environment variables.
	apiToken := "EAATZAxfQTVYQBAD5RIvKCpLEK5BQ4TF7V2l6S4OYcWHZAxZAwQ1va2x5zGNZAgEke8ZC7Mik8CKOcwqPmSLZBrZB2PzBaXEeOvhvoxfHwjelZBMLZCGZCOvflQJ1cCSH2nPfdOVih79WoQK0F47I5BI6wetibxz0eTlsiWFv9gPbllZBgZDZD"
	os.Setenv("TOKEN", apiToken)

	cbMessenger.MessageReceived = MessageReceived

	// API endpoints
	http.HandleFunc("/webhook", cbMessenger.Handler)
	http.HandleFunc("/campaign", campaignHandler)

	// Heroku has its own env PORT. If not available use 3001 (for local development)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
