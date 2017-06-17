package campaign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Content struct {
	Text string `json:"text"`
}

type UserAction struct {
	NodeType string `json:"type"`
	Target   string `json:"target"`
	Label    string `json:"label"`
}

type CampaignNode struct {
	Effect      string `json:"effect"`
	Content     `json:"content"`
	UserActions []UserAction `json:"user_actions"`
}

type Campaign struct {
	PublishAt 	string   `json:"publish_at"`
	Topics    	[]string `json:"topics"`
	Id 			string   `json:"uuid"`
	Name      	string   `json:"name"`
	RootNode  	string   `json:"root_node"`
	Nodes     	map[string]CampaignNode{}
}

// TODO: will need to remove this later once the campaign API can send encoded JSON 
func init() {
	raw, err := ioutil.ReadFile("../../../campaign-node.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var campaign Campaign
	json.Unmarshal(raw, &campaigns)

	fmt.Println(toJson(campaign))

	// Encode the Campaign Struct to JSON
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(campaign)
	res, err := http.Post("/campaign", "application/json; charset=utf-8", b)
}
