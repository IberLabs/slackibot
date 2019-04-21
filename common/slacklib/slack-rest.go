package slacklib

import (
	"github.com/nlopes/slack"
	"fmt"
	"errors"
)

// Get channel list using REST API
func GetUserList(token string) ([]slack.User, error) {
	api := slack.New(token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return users, err
}


/**
	Send direct message to a user for the first time
 */
func SendNewDirectMessage(token string, msg string, user string) (error, string) {
	api := slack.New(token)
	var err error = nil
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	noOp, _, channelID, _ := api.OpenIMChannel(user)
	if channelID == "" || noOp == false  {
		err = errors.New("Error sending direct message to " + user)
	}

	return err, channelID
}


/**
	Slack real-time API (websockets) does not provide chanel name. It has to be compared with REST API channel list
	@channelId	string				Channel ID
	@channels	[]slack.Channel		Array of channels gathered from REST API
 */
func GetChannelRealName(channelID string, channels []slack.Channel) (string, bool) {
	isIM := false

	if len(channelID) > 0 && channelID[0] == 'D' {
		channelID = "[DIRECT-MSG]"
		isIM = true
	}else {
		for _, v := range channels {
			if v.ID == channelID {
				return v.NameNormalized, isIM
			}
		}
	}

	return channelID, isIM
}


/**
	Slack RT API does not provide real name. It has to be compared with REST API user list
 */
func GetUserRealName(userID string, users []slack.User) (string) {
	for _, v := range users {
		if v.ID == userID {
			return v.RealName
		}
	}

	return userID
}
