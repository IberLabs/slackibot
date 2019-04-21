package slacklib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
	"github.com/nlopes/slack"
)

type Attachment = slack.Attachment

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.

type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

// Get channel list using REST API
func GetChannelList(token string) ([]slack.Channel, error) {
	api := slack.New(token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	//for _, chn := range channels {
	//	fmt.Printf("ID: %s, Name: %s\n", chn.ID, chn.NameNormalized )
	//}

	return channels, err
}


/**
	Send direct IM message to a user.
	Specifying channelID means that this conversation was opened before.
 */
func SendDirectMessage(token string, msg string, user string, channelID string, attachment slack.Attachment) (error, string) {
	api := slack.New(token)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)
	var err error

	// No channelID means that no previous direct conversation created with the user
	if channelID == "" {
		err, channelID = SendNewDirectMessage(token, msg, user)
	}

	if err == nil {

		if attachment.Text != "" {
			_, _, err = api.PostMessage(channelID, slack.MsgOptionText(msg, false),  slack.MsgOptionAttachments(attachment))
		}else{
			_, _, err = api.PostMessage(channelID, slack.MsgOptionText(msg, false))
		}
	}

	return err, channelID
}

// Websocket real time API to wait for events
// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func slackStart(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

func GetMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

var counter uint64

func PostMessage(ws *websocket.Conn, m Message) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func SlackConnect(token string) (*websocket.Conn, string) {
	wsurl, id, err := slackStart(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	return ws, id
}
