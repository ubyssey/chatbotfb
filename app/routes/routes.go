package routes

import (
	"fmt"
	"net/http"

	"github.com/ubyssey/chatbotfb/app/lib/chatbot"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/webhook", chatbot.CbMessenger.Handler)
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is working!~!@$!@$!@$!@$!@$!$!@")
	})

	return router
}
