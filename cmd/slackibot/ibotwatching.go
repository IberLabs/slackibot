package main

import (
	"slackibot/assets"
	c "slackibot/common"
	"fmt"
	"github.com/nlopes/slack"
	"golang.org/x/net/websocket"
	"slackibot/common/slacklib"
	"strings"
	"time"
)


/**
	Detect if the bot has been mentioned (interacti)
 */
func mentionedMe(ws *websocket.Conn, m slacklib.Message) {
	// if so try to parse if
	parts := strings.Fields(m.Text)
	if len(parts) == 3 && parts[1] == "info" {
		// looks good, get the quote and reply with the result
		go func(m slacklib.Message) {
			m.Text = "iBot - Immediate alerts and efficient support.\n"
			slacklib.PostMessage(ws, m)
		}(m)
		// NOTE: the Message object is copied, this is intentional
	} else {
		// huh?
		m.Text = fmt.Sprintf("sorry, I can't understand your request. Check \"info\" command\n")
		slacklib.PostMessage(ws, m)
	}
}


/**
	Check incoming messages and look for special words that trigger alerts
 */
func findAlertWords(ws *websocket.Conn, m slacklib.Message, configuration c.Configuration, channels[] slack.Channel) bool {
	parts := strings.Fields(m.Text)
	countParts := len(parts)
	alert := false
	alertWord := ""

	for i := 0; i < countParts; i++ {
		for x := 0; x < len(assets.AlertWords); x++ {
			if strings.ToLower(parts[i]) == assets.AlertWords[x] {
				alertWord = parts[i]
				logMsg := "Special word detected (word: " + alertWord + " channel: " + m.Channel + ", user: " + m.User + " )\n"
				c.Display(logMsg, false, true)
				alert = true
				alertWord = parts[i]
				break
			}
		}
		if alert {
			break
		}
	}

	// IM will be true in case of direct message
	// MODIFIED FOR POC
	channelRealName, _ := slacklib.GetChannelRealName(m.Channel, channels)
	if alert && !c.GetAlertFromDBLastMins(channelRealName, ALERTS_TIMEWAIT) {
		botResponse(ws, channelRealName, alertWord, m, configuration, false)
	}

	return alert
}


/**
	Trigger a bot response
	bool IM = true if direct message instead of channel
 */
func botResponse(ws *websocket.Conn, channelRealName string, alertWord string, m slacklib.Message, configuration c.Configuration, IM bool){
	msgResponse := "Special word '" + alertWord + "' detected, triggering alert. A team member will bring you support in short.\n"

	if isOutOfOffice() {
		msgResponse = msgResponse + "Out of office hours. Activating oncall support.\n"
	}

	// Standard channel response message
	slacklib.PostMessage(
		ws,
		slacklib.Message{m.Id, m.Type,m.Channel, msgResponse, ""})

	// Send alert if "office-hours" or almost register it into log DB
	c.EvalSendAlert(configuration.TokenOpsgenie, m.Text, channelRealName, m.User, alertWord, false)
}


/**
	Check if current time matches with out-of-office time-window
 */
func isOutOfOffice() bool{
	result := false

	thisWeekDay := time.Now().Weekday()
	thisHour := time.Now().Hour()
	if thisWeekDay == TIMETABLE_SATURDAY || thisWeekDay == TIMETABLE_SUNDAY {
		// Non working day
		result = true
	}

	if !result && (thisHour < TIMETABLE_WATCHING_START || thisHour > TIMETABLE_WATCHING_END -1) {
		// No office hours
		result = true
	}

	dateString := time.Now().Format("01-02-2016 15:04:05")
	if result {
		dateString += " Out-of-office hours"
	}else{
		dateString += " Office hours"
	}
	c.Display(time.Now().Format("01-02-2016 15:04:05") + dateString, false, true)

	return result
}