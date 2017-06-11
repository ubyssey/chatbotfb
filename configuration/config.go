package configuration

import (
	"error"
	"os"

	"github.com/ubyssey/chatbotfb/app/utils/jsonparser"
	"github.com/ubyssey/chatbotfb/app/utils/printlogger"
)

type Configuration struct {
	Database `json:"Database"`
}

type Database struct {
	MongoDB `json:"MongoDB"`
}

type MongoDB struct {
	LocalPort string `json:"LocalPort"`
	LocalURL  string `json:"LocalURL"`
	Name      string `json:"Name"`
}

var Config = &Configuration{}

func init() {
	raw, err := ioutil.ReadFile("./config.json")

	if err != nil {
		printlogger.Log(err.Error())
		os.Exit(1)
	}

	parseError := jsonparser.Parse(raw, Config)

	if parseError != nil {
		printLogger.Log(parseError.Error())
	}
}
