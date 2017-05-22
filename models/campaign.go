package models

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// )

// type Content struct {
// 	Text string `json:"text"`
// }

// type UserAction struct {
// 	NodeType string `json:"type"`
// 	Target   string `json:"target"`
// 	Label    string `json:"label"`
// }

// type CampaignNode struct {
// 	Effect      string `json:"effect"`
// 	Content     `json:"content"`
// 	UserActions []UserAction `json:"user_actions"`
// }

// type Campaign struct {
// 	PublishAt string   `json:"publish_at"`
// 	Topics    []string `json:"topics"`
// 	Name      string   `json:"name"`
// 	RootNode  string   `json:"root_node"`
// 	Nodes     map[string]CampaignNode
// }

// func getCampaign() {
// 	return
// }

// // Test campaign
// func StartTestCampaign() {
// 	// TODO: change to root directory since this is being called from elsewhere
// 	raw, err := ioutil.ReadFile("./campaign-node.json")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	var campaigns []Campaign
// 	json.Unmarshal(raw, &campaigns)

// 	for _, c := range campaigns {
// 		fmt.Println(toJson(c))
// 	}

// 	fmt.Println(toJson(campaigns))
// }

// func toJson(c interface{}) string {
// 	bytes, err := json.Marshal(c)

// 	if err != nil {
// 		fmt.Println(err.Error())
// 		os.Exit(1)
// 	}

// 	return string(bytes)
// }
