package messageoptions

func ShowMenu(senderID string) {
	mq := messenger.MessageQuery{}
	mq.RecipientID(senderID)

	menuButtonOptions := getMenuButtonOptions()
	mq.Template(
		template.GenericTemplate{
			Title:    "Menu",
			ImageURL: "https://upload.wikimedia.org/wikipedia/commons/thumb/b/b1/Hot_dog_with_mustard.png/1200px-Hot_dog_with_mustard.png",
			Buttons:  menuButtonOptions,
		},
	)
}

// Start a new campaign for the user
func StartCampaign() {
	startCampaign := campaign.Campaign{}
	startCampaignMissingErr := campaignCollection.FindId(startCampaignId).One(&startCampaign)

	if startCampaignMissingErr != nil {
		chatbot.DefaultMessage(senderID, "A start campaign was not found.")
		return
	}

	// Get the button templates to be shown to the user
	buttonsOptions, buttonOptionsErr := chatbot.GetButtonTemplateOptions(
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

func getMenuButtonOptions() []template.Buttons {

}
