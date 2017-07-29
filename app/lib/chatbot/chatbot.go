package chatbot

import (
	"os"

	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	// Messenger SDK
	CbMessenger = &messenger.Messenger{
		AccessToken: os.Getenv("TOKEN"),
		VerifyToken: os.Getenv("TOKEN"),
	}
	respText string
)

func DefaultMessage(recipient string, messages ...interface{}) {
	if len(messages) > 0 {
		respText = messages
	} else {
		respText = "Sorry, I don't understand your message."
	}

	resp, msgErr := CbMessenger.SendSimpleMessage(
		recipient,
		fmt.Sprintf(respText),
	)

	if msgErr != nil {
		printlogger.Log(msgErr.Error())
	}
	printlogger.Log("%+v", resp)
}

func GetCampaignButtonTemplateOptions(campaignId string, userActions []campaign.UserAction) ([]template.Button, error) {
	// A button slice to hold each button option to be shown to the user
	buttonSlice := []template.Button{}
	var button template.Button

	for _, currUserAction := range userActions {
		// Reset the button struct every loop
		button = template.Button{}

		// If the node type is a "link", create a NewWebURLButton template. Otherwise,
		// if it is a "node", then create a postback Button template with its payload
		actionNodeType := currUserAction.NodeType

		if actionNodeType == "link" {
			button = template.NewWebURLButton(
				currUserAction.Label,
				currUserAction.Target,
			)
		} else if actionNodeType == "node" {
			payloadOption := payload.Payload{
				CampaignId: campaignId,
				Event: &user.Event{
					NodeType: actionNodeType,
					Target:   currUserAction.Target,
					Label:    currUserAction.Label,
				},
			}

			payloadOptionString, payloadOptionParsingErr := jsonparser.ToJsonString(payloadOption)

			if payloadOptionParsingErr != nil {
				return nil, payloadOptionParsingErr
			}

			button = template.Button{
				Type:    template.ButtonTypePostback,
				Payload: payloadOptionString,
				Title:   currUserAction.Label,
			}
		}

		// If the button struct is not empty, append it to the buttonSlice
		if button.Type != "" && button.Title != "" {
			buttonSlice = append(buttonSlice, button)
		}
	}
	return buttonSlice, nil
}
