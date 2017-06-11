package routes

import (
	"github.com/gorilla/mux"
)

func routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/webhook")
}
