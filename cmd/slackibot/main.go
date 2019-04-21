package main

import (
	"flag"
	"fmt"
	"github.com/google/logger"
	"github.com/nlopes/slack"
	"golang.org/x/net/websocket"
	"os"
	"slackibot/assets"
	c "slackibot/common"
	"slackibot/common/slacklib"
	"strings"
	"time"
)


const GLBL_VERSION					= "0.1.3"
const GLBL_VERBOSE					= false
const SLCK_EVENT_TYPING				= "user_typing"
const SLCK_EVENT_MESSAGE			= "message"
const ALERTS_TIMEWAIT				= 5
const TIMETABLE_WATCHING_START		= 9					// Watching period start hour
const TIMETABLE_WATCHING_END		= 18				// Watching period end hour
const TIMETABLE_SATURDAY			= 6					// Saturday day of week
const TIMETABLE_SUNDAY				= 0					// Sunday day of week

const logPath 						= "./ibot.log"
const configFile 					= "./config.json"
const dbFile						= "./bot.db"


/**
	Main application entry point.
 */
func main() {
	// Cobra commands (will seek within "common" folder for commands)
	if !c.Execute(){
		// Stop execution. Application command invoked.
		os.Exit(1)
	}

	// Check and init local DB
	c.InitDBcheck(dbFile)

	// Init global logging
	flag.Parse()
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		c.Display("Failed to open log file: " +  string(err.Error()), true, true)
	}
	defer lf.Close()

	if(GLBL_VERBOSE) {
		var verboseMode= flag.Bool("verbose", false, "print info level logs to stdout")
		defer logger.Init("slackibot", *verboseMode, true, lf).Close()
	}else{
		defer logger.Init("slackibot", false, true, lf).Close()
	}
	c.Display("slackibot started!", false, true)

	// Load general configuration from json file
	configuration := c.LoadConfig(configFile)

	// Init channel list
	channels, err := slacklib.GetChannelList(configuration.TokenSlack)
	if err != nil {
		c.Display( "Error retrieving slack channel list", false, true )
		os.Exit(1)
	}

	// Init user list
	users, err := slacklib.GetUserList(configuration.TokenSlack)
	if err != nil {
		c.Display( "Error retrieving slack user list", false, true )
		os.Exit(1)
	}

	// Enter the watching bucle (wait for events)
	slackBotWatching(configuration, channels, users)
}


/**
	slackBotWatching: Watching bucle. Wait for slack websocket events.
 */
func slackBotWatching(configuration c.Configuration, channels []slack.Channel, users []slack.User){
	// start a websocket-based Real Time API session
	ws, id := slacklib.SlackConnect(configuration.TokenSlack)
	c.Display( "slackibot ready, press ^C to exit", false, true )

	for {
		// read each incoming message
		m, err := slacklib.GetMessage(ws)
		if err != nil {
			c.Display(err.Error(), true, true)
			continue
		}

		// cName is channel name. IM is true if it is a direct message.
		chName, IM := slacklib.GetChannelRealName(m.Channel, channels)
		// uName is the user name.
		uName := slacklib.GetUserRealName(m.User, users)

		// If channel starts with letter 'D' it is a direct IM message
		c.Display( "WebSocket event: " + m.Type + " chan: " +  chName + " (" + m.Channel + ")  usr: " + uName, false, true)

		// When uName is blank means that the message is mine.
		if m.Type == "message" && uName != "" {

			if IM{
				// Direct message conversation
				// NO DIRECT IM CONVERSATION ON slackibot LITE VERSION
			}else{
				// Public conversation
				findAlertWords(ws, m, configuration, channels)

				// Detect if bot has been mentioned
				if strings.HasPrefix(m.Text, "<@"+id+">") {
					mentionedMe(ws, m)
					c.Display("Event mentionedMe :: "+m.Text, false, true)
				}
			}

		}else{
			// Other event types.
		}
	}
}


/**
	Detect if the bot has been mentioned
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
		msgResponse = msgResponse + "Out of office hours.\nActivating oncall support.\n"
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

	if !result && (thisHour >= TIMETABLE_WATCHING_START || thisHour <= TIMETABLE_WATCHING_END -1) {
		// No office hours
		result = true
	}

	if result {
		c.Display(time.Now().Format("01-02-2016 15:04:05") + " Out of office hours response", false, true)
	}

	return result
}