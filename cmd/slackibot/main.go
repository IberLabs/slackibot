package main

import (
	"flag"
	"github.com/google/logger"
	"github.com/nlopes/slack"
	"os"
	c "slackibot/common"
	"slackibot/common/slacklib"
	"strings"
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
	Application main entry point.
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

	var verboseMode= flag.Bool("verbose", GLBL_VERBOSE, "print info level logs to stdout")
	defer logger.Init("slackibot", *verboseMode, true, lf).Close()
	c.Display("slackibot " + GLBL_VERSION + " started!", false, true)

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