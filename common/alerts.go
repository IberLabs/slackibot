package common


import (
	"github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	"fmt"
	ogcli "github.com/opsgenie/opsgenie-go-sdk/client"
)

/**
	Evaluate if an alert has to be sent.
 */
func EvalSendAlert(OpsgenieToken string, txt string, channel string, user string, words string, noTrigger bool){

	prefix := ""

	if noTrigger {
		prefix = "[SILENT] "
	}

	// Register event into DB
	AddAlertToDB(user, channel, words, prefix + txt)

	if !noTrigger {
		// Trigger "special word" alert if time is within office hours
		triggerAlert(OpsgenieToken, txt, channel, user)
	}

}

/**
	Send and opsgenie Alert
 */
func triggerAlert(OpsGenieToken string, txt string, channel string, user string) {
	cli := new(ogcli.OpsGenieClient)
	cli.SetAPIKey(OpsGenieToken)

	alertCli, _ := cli.AlertV2()

	teams := []alertsv2.TeamRecipient{
		&alertsv2.Team{Name: "CI_Tools_Operations"},
		&alertsv2.Team{ID: "CI_Tools_Operations"},
	}

	/**
	visibleTo := []alertsv2.Recipient{
		&alertsv2.Team{ID: " "},
		&alertsv2.Team{Name: " "},
		&alertsv2.User{ID: " "},
		&alertsv2.User{Username: " "},
	}
	*/

	request := alertsv2.CreateAlertRequest{
		Message:     "Slack channel \"" + channel + "\" alert. Read channel and bring support.",
		Alias:       "slackAlert",
		Description: user + " said: " + txt,
		Teams:       teams,
		/* VisibleTo:   visibleTo, */
		Actions:     []string{"action1", "action2"},
		Tags:        []string{"tag1", "tag2"},
		Details: map[string]string{
			"channel":  channel,
			"key2": "value2",
		},
		Entity:   "entity",
		Source:   "source",
		Priority: alertsv2.P2,
		/* User:     "user@opsgenie.com", */
		Note:     " ",
	}

	response, err := alertCli.Create(request)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Create request ID: " + response.RequestID)
	}
}
