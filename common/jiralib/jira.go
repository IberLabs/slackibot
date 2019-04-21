package jiralib

import (
	"net/http"
	"fmt"
	goJira "github.com/andygrunwald/go-jira"
	"slackibot/common/jiralib/oauth"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
	"slackibot/common"
)

type EnumAuthType int
const (
	EnumOAuth 	EnumAuthType = 1
	EnumCookie 	EnumAuthType = 2
)

type conditionMatch = ConditionMatch

type ConditionGroup struct {
	ConditionMatch []conditionMatch
	Match          string
}

type JiraConfig struct {
	Login 			string
	Pw				string
	BaseUrl			string
	ApiUrl			string
	SessUrl			string
}

const api_issue_endpoint = "issue"

func CreateJiraIssue(jiraCfg JiraConfig) {

	// TODO: Allow Oauth authentication.
	jiraClient := getJiraClientCookies(EnumCookie, jiraCfg)

	m := MessageCreateIssue{
		Fields : MCIFields{
			Project : 	MCIProject{
				Key : 	"OCI",
			},
			Summary : 	"This is a test issue",
			Description:"This is a test description",
			IssueType: MCIIssueType{
				Name: 	"Incident",
			},
			Impact: MCIImpact{
				Value: "High",
			},
			Urgency: MCIUrgency{
				Value: "High",
			},
			Components: [] MCIComponents{
				{ Name: "Jenkins"},
			},
			GOA_Environment: [] MCIEnvironment{
				{ Value: "Production"},
			},
		},
	}

	mJson, err 	:= json.Marshal(m)
	PanicIfError(err, "")
	jsonStr 	:= []byte(mJson)

common.Display(string(jsonStr), false, true)

	resp, body 	:= postRequest(
		jiraCfg.BaseUrl + jiraCfg.ApiUrl + api_issue_endpoint,
		jiraCfg.Login,
		jiraCfg.Pw,
		jiraClient,
		jsonStr)

	common.Display(resp.Status, false, true )
	common.Display(body, false, true)

}


func postRequest(apiUrl string, user string, pw string, jiraCLient *goJira.Client, msg []byte) (*goJira.Response, string) {
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, pw)

	//fmt.Println("REQUEST-BODY:", string(msg))
	resp, err := jiraCLient.Do(req, nil)
	body, _ := ioutil.ReadAll(resp.Body)
	PanicIfError(err, string(body[:]))

	return resp, string(body)
}


func createMessage(start uint, filters Filters) []byte {
	str := createJqlRequest(filters)
	fmt.Println(str)
	m := JqlMessage{StartAt: start, MaxResults: 1000, Fields: filters.Fields, Jql: str}
	mJson, err := json.Marshal(m)
	PanicIfError(err, "")
	jsonStr := []byte(mJson)
	return jsonStr
}


func createJqlRequest(filters Filters) string {
	cmArray := make([]conditionMatch, len(filters.Project), len(filters.Project))
	cmArray[0] = conditionMatch{Key: "project", Match: " in ", Value: filters.Project}
	now := time.Now()
	cgArray := make([]string, 2, 2)
	cmArray2 := make([]conditionMatch, len(filters.Conditions), len(filters.Conditions))
	for i, filter := range filters.Conditions {
		if filter.Type == "time" {
			value, _ := strconv.Atoi(filter.Value)
			then := now.AddDate(0, 0, -value)
			cmArray2[i] = conditionMatch{Key: filter.Key, Match: filter.Match, Value: then.Format("2006-01-02")}
		} else {
			cmArray2[i] = conditionMatch{Key: filter.Key, Match: filter.Match, Value: filter.Value}
		}
	}
	cGroup1 := ConditionGroup{ConditionMatch: cmArray, Match: "OR"}
	cGroup2 := ConditionGroup{ConditionMatch: cmArray2, Match: "AND"}
	cgArray = append(cgArray, CreateConditionGroup(cGroup1))
	cgArray = append(cgArray, CreateConditionGroup(cGroup2))
	str := GroupConditionGroup(cgArray, "AND")
	return str
}


func CreateConditionGroup(conditionG ConditionGroup) string {
	var str string
	for _, conditionM := range conditionG.ConditionMatch {
		if len(conditionM.Match) > 0 {
			if len(str) > 0 {
				str = str + " " + conditionG.Match + " " + CreateJql(conditionM)
			} else {
				str = CreateJql(conditionM)
			}
		}
	}
	return str
}

func GroupConditionGroup(conditionG []string, match string) string {
	var str string
	for _, condition := range conditionG {
		if len(condition) > 0 {
			if len(str) > 0 {
				str = str + match + "(" + condition + ")"
			} else {
				str = "(" + condition + ")"
			}
		}
	}
	return str
}


func CreateJql(conditionM conditionMatch) string {
	str := string(conditionM.Key + conditionM.Match + conditionM.Value)
	return str
}


func getJiraClientCookies(authType EnumAuthType, cfg JiraConfig) * goJira.Client {
	jiraClient, err := goJira.NewClient(nil, cfg.BaseUrl)
	PanicIfError(err, "")
	res, err := jiraClient.Authentication.AcquireSessionCookie(cfg.Login, cfg.Pw)
	if err != nil || res == false {
		fmt.Printf("Result: %+v\n", res)
		panic(err)
	}
	return jiraClient

}


func getJiraClientOauth(authType EnumAuthType, cfg JiraConfig) *http.Client {
	jiraClient := oauth.GetJIRAClient()

	return jiraClient
}


func PanicIfError(err error, details string) {
	if err != nil {
		fmt.Println(details)
		panic(err)
	}
}
