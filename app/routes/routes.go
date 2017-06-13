package routes

import (
	"github.com/ubyssey/chatbotfb/app/controllers"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/webhook", chatbot.CbMessenger.Handler)
	router.HandleFunc("/campaign", controllers.GetCampaign)

	return router
}
