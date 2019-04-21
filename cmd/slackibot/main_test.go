package main

import (
	"testing"
	jiraLib "slackibot/common/jiralib"
	"slackibot/common"
	"os"
)

const testLogPath 					= "./bin/ibot.log"
const testConfigFile 				= "./bin/config.json"
const testDbFile					= "./bin/bot.db"

var configuration common.Configuration

// Initializer / setup fixtures
func TestMain(m *testing.M) {
	loadFixtures()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestCreateJiraIssue(t *testing.T) {
	loadFixtures()

	jiraLib.CreateJiraIssue(
		jiraLib.JiraConfig {
			Login : configuration.Jira.CookieLogin.Username,
			Pw: 	configuration.Jira.CookieLogin.Password,
			BaseUrl:configuration.Jira.BaseUrl,
			ApiUrl: configuration.Jira.ApiUrl,
			SessUrl:configuration.Jira.CookieLogin.BaseSessUrl,
		},
	)

}

func loadFixtures(){
	configuration = loadConfig(testConfigFile)
}

func tearDown() {
	// Nothing to do right now
}


