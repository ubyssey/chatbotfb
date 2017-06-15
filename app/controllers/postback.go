package controllers

import (
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

func init() {
	chatbot.CbMessenger.Postback = postPostback
}

func postPostback(event messenger.Event, opts messenger.MessageOpts, pb messenger.Postback) {

}
