package messageactions

import (
	"github.com/ubyssey/chatbotfb/app/database"
	"github.com/ubyssey/chatbotfb/app/external/ubysseydotca"
	"github.com/ubyssey/chatbotfb/app/lib/chatbot"
	"github.com/ubyssey/chatbotfb/app/models/campaign"
	"github.com/ubyssey/chatbotfb/app/server/payload"
	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
	"github.com/ubyssey/chatbotfb/configuration"

	"github.com/maciekmm/messenger-platform-go-sdk/template"
	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

var (
	startCampaignId = "6z479nb9-3x2f-23gs-g2dz-abc10625xc68"
)

// Send the user actions for the next campaign node
func SendNextCampaignNodeActions(senderID string, payloadStruct payload.Payload) {
	dbName := configuration.Config.Database.MongoDB.Name
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")
	currentCampaign, currentCampaignErr := campaign.GetCampaignStruct(campaignCollection, payloadStruct.CampaignId)
	if currentCampaignErr != nil {
		printlogger.Log(currentCampaignErr.Error())
		printlogger.Log("Error finding the campaign :%s", payloadStruct.CampaignId)
		return
	}

	if campaignNode, ok := currentCampaign.Nodes[payloadStruct.Event.Target]; ok {
		// If a node still has children, send a message with those children node options,
		// otherwise send the final message of the current campaign
		if len(campaignNode.UserActions) > 0 {
			mq := messenger.MessageQuery{}
			mq.RecipientID(senderID)

			buttonsOptions, buttonOptionsErr := chatbot.GetCampaignButtonTemplateOptions(
				payloadStruct.CampaignId,
				campaignNode.UserActions,
			)

			if buttonOptionsErr != nil {
				printlogger.Log(buttonOptionsErr.Error())
				return
			}

			// Generic Message Query template to be sent to the user
			mq.Template(
				template.GenericTemplate{
					Title:   currentCampaign.Name,
					Buttons: buttonsOptions,
				},
			)

			resp, msgErr := chatbot.CbMessenger.SendMessage(mq)

			if msgErr != nil {
				printlogger.Log(msgErr.Error())
			}

			printlogger.Log("%+v", resp)
		} else {
			chatbot.DefaultMessage(senderID, "")
		}
	} else {
		printlogger.Log("Campaign Node target %s not found for user %s", payloadStruct.Event.Target, senderID)
	}
}

// Shows the main menu to the user
func ShowMenu(senderID string) {
	menuButtonOptions, menuButtonOptionsErr := getMenuButtonOptions()

	if menuButtonOptionsErr != nil {
		printlogger.Log(menuButtonOptionsErr.Error())
		return
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)
	mq.Template(
		template.GenericTemplate{
			Title:    "Menu",
			ImageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b1/Hot_dog_with_mustard.png/1200px-Hot_dog_with_mustard.png",
			Buttons:  menuButtonOptions,
		},
	)

	resp, msgErr := chatbot.CbMessenger.SendMessage(mq)
	if msgErr != nil {
		printlogger.Log(msgErr.Error())
	}
	printlogger.Log("%+v", resp)
}

// Shows the list template for the menu postback (a ubyssey.com link)
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/list-template
func ShowMenuListTemplate(senderID string, payloadStruct payload.Payload) {
	articles, articlesErr := ubysseydotca.GetHttpRequestApi(payloadStruct.UrlLink)
	if articlesErr != nil {
		printlogger.Log(articlesErr.Error())
		return
	}

	listElementSlice := []template.ListElement{}

	for _, article := range articles.Results {
		listElement := template.ListElement{
			Title:    article.Headline,
			ImageURL: article.FeaturedImage.Url,
			Subtitle: article.Snippet,
			DefaultAction: template.DefaultAction{
				Type:                "web_url",
				URL:                 article.Url,
				MessengerExtensions: true,
				WebviewHeightRatio:  "tall",
			},
		}
		listElementSlice = append(listElementSlice, listElement)
	}

	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	mq.Template(
		template.ListTemplate{
			Elements: listElementSlice,
		},
	)

	resp, msgErr := chatbot.CbMessenger.SendMessage(mq)
	if msgErr != nil {
		printlogger.Log(msgErr.Error())
	}
	printlogger.Log("%+v", resp)
}

// Start a new campaign for the user
func StartCampaign(senderID string) {
	dbName := configuration.Config.Database.MongoDB.Name
	campaignCollection := database.MongoSession.DB(dbName).C("campaigns")

	startCampaign := campaign.Campaign{}
	startCampaignMissingErr := campaignCollection.FindId(startCampaignId).One(&startCampaign)

	if startCampaignMissingErr != nil {
		chatbot.DefaultMessage(senderID, "A start campaign was not found.")
		return
	}

	// Get the button templates to be shown to the user
	buttonsOptions, buttonOptionsErr := chatbot.GetCampaignButtonTemplateOptions(
		startCampaign.UUID,
		startCampaign.Nodes[startCampaign.RootNode].UserActions,
	)

	if buttonOptionsErr != nil {
		printlogger.Log(buttonOptionsErr.Error())
		return
	}

	// Initialize a message query
	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	mq.Template(
		template.GenericTemplate{
			Title:   startCampaign.Name,
			Buttons: buttonsOptions,
		},
	)

	resp, msgErr := chatbot.CbMessenger.SendMessage(mq)
	if msgErr != nil {
		printlogger.Log(msgErr.Error())
	}
	printlogger.Log("%+v", resp)
}

// Get the button options for the menu
func getMenuButtonOptions() ([]template.Button, error) {
	// A button slice to hold each button option to be shown to the user
	buttonSlice := []template.Button{}

	// Set the button titles in the menu
	menuOptionsMap := map[string]string{}
	menuOptionsMap["Top stories"] = "https://www.ubyssey.ca/api/articles/?limit=4"
	menuOptionsMap["News"] = "https://www.ubyssey.ca/api/articles/?section=1&limit=4"
	menuOptionsMap["Sports"] = "https://www.ubyssey.ca/api/articles/?section=2&limit=4"
	menuOptionsMap["Culture "] = "https://www.ubyssey.ca/api/articles/?section=3&limit=4"

	// For each menu option, create a button and add it to the button slice
	for buttonTitle, buttonUrl := range menuOptionsMap {
		payloadStruct := payload.Payload{
			UrlLink: buttonUrl,
		}

		payloadString, payloadParsingErr := jsonparser.ToJsonString(payloadStruct)

		if payloadParsingErr != nil {
			return nil, payloadParsingErr
		}

		button := template.Button{
			Type:    template.ButtonTypePostback,
			Payload: payloadString,
			Title:   buttonTitle,
		}

		buttonSlice = append(buttonSlice, button)
	}

	return buttonSlice, nil
}
