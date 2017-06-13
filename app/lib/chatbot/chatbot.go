package chatbot

import (
	"os"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	// Messenger SDK
	CbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
)
