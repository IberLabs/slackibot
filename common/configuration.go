package common

import (
	"github.com/tkanos/gonfig"
	"os"
)

/**
	Description of configuration file that should be loaded from config.json file in the bin folder
 */

type Configuration struct {
	TokenSlack   	string
	TokenOpsgenie 	string
	TokenJira		string
	Jira struct {
		BaseUrl		string
		ApiUrl		string
		SearchUrl	string
		ProjectUrl	string
		WorkflowUrl	string
		CookieLogin struct {
			Username		string
			Password 		string
			BaseSessUrl    	string
		}
		Token		string
	}
}

/**
	Load/Initialize application configuration
 */
func LoadConfig(confFile string) Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf(confFile, &configuration)
	if err != nil {
		Display(err.Error(), true, true)
		os.Exit(1)
	}
	return configuration
}
