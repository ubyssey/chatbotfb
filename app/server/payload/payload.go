package payload

import (
	"github.com/ubyssey/chatbotfb/app/models/user"
	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"

	"gopkg.in/maciekmm/messenger-platform-go-sdk.v4"
)

type Payload struct {
	IsUrlLink  bool        `json:"is_url_link,omitempty"`
	UrlLink    string      `json:"url_link,omitempty"`
	CampaignId string      `json:"campaign_id,omitempty"`
	Event      *user.Event `json:"event,omitempty"`
}

func GetPayloadStruct(pb messenger.Postback) (Payload, error) {
	// Get the payload from the postback message
	payloadStruct := Payload{}
	paylodStructParseErr := jsonparser.Parse([]byte(pb.Payload), &payloadStruct)
	return payloadStruct, paylodStructParseErr
}
