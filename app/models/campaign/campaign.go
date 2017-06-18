package campaign

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

type Content struct {
	Text string `json:"text" bson:"text"`
}

type UserAction struct {
	NodeType string `json:"type" bson:"nodeType"`
	Target   string `json:"target" bson:"target"`
	Label    string `json:"label" bson:"label"`
}

type CampaignNode struct {
	Effect      string `json:"effect" bson:"effect"`
	Content     `json:"content" bson:"content"`
	UserActions []UserAction `json:"user_actions" bson:"userActions"`
}

type Campaign struct {
	PublishAt 	string   `json:"publish_at" bson:"publish_at"`
	Topics    	[]string `json:"topics" bson:"topics"`
	UUID 		string   `json:"uuid" bson:"_id,omitempty"`
	VersionUUID string	 `json:"version_uuid" bson:"versionUUID"`
	Name      	string   `json:"name" bson:"name"`
	RootNode  	string   `json:"root_node" bson:"rootNode"`
	Nodes     	map[string]CampaignNode{} `json:"nodes" bson:"nodes"`
}