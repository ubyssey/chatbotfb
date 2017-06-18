package payload

import (
	"github.com/ubyssey/chatbotfb/app/models/user"
)

type Payload struct {
	CampaignId string      `json:"campaign_id"`
	Event      *user.Event `json:"event"`
}
