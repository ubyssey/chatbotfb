package routes

import (
	"github.com/ubyssey/chatbotfb/app/controllers"

	"github.com/gorilla/mux"
	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	// Messenger SDK
	CbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	CbMessenger.MessageReceived = controllers.GetMessage

	router.HandleFunc("/webhook", CbMessenger.Handler)
	router.HandleFunc("/campaign", controllers.GetCampaign)

	return router
}
