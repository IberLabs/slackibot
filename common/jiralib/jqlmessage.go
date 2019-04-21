package jiralib

type JqlMessage struct {
	Jql        string   `json:"jql"`
	StartAt    uint     `json:"startAt"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields"`
	//FieldsByKeys bool     `json:"fieldsByKeys"`
}
