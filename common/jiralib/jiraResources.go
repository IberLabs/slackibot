package jiralib

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
)

type Properties struct {
	BaseUrl     string `json:"baseUrl"`
	SearchUrl   string `json:"searchUrl"`
	ProjectUrl  string `json:"projectUrl"`
	WorkflowUrl string `json:"workflowUrl"`
}

type Filters struct {
	Project    string           `json:"project"`
	Fields     []string         `json:"fields"`
	Conditions []ConditionMatch `json:"conditions"`
}

type ConditionMatch struct {
	Key   string `json:"key"`
	Match string `json:"match"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Token struct {
	Token string `json:"token"`
}

//----------------------------------------------------------------------------------

type MessageCreateIssue struct {
	Fields 				MCIFields		`json:"fields"`
}

type MCIFields struct{
	Project 				MCIProject		`json:"project"`
	Summary 				string			`json:"summary"`
	Description 			string			`json:"description"`
	IssueType 				MCIIssueType	`json:"issuetype"`
	Impact					MCIImpact		`json:"customfield_10203"`
	Urgency 				MCIUrgency		`json:"customfield_10204"`
	Components			[] MCIComponents	`json:"components"`
	GOA_Environment	 	[] MCIEnvironment	`json:"customfield_10205"`
}

type MCIComponents struct {
	Name					string			`json:"name"`
}

type MCIUrgency struct {
	Value					string			`json:"value"`
}

type MCIImpact struct {
	Value					string			`json:"value"`
}

type MCIEnvironment struct {
	Value					string			`json:"value"`
}

type MCIProject struct {
	Key					string				`json:"key"`
}

type MCIIssueType struct {
	Name 				string				`json:"name"`
}

//----------------------------------------------------------------------------------

func GetProperties(filePath string) Properties {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Properties
	json.Unmarshal(raw, &c)
	return c
}

func GetFilters(filename *string) Filters {
	raw, err := ioutil.ReadFile("./resources/" + *filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Filters
	json.Unmarshal(raw, &c)
	return c
}
