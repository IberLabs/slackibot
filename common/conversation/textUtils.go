package conversation

import "strings"

type AlertRequestItems struct {
	err				error
	tool			Tool
	origText		string
	alertText		string
	tags		[]	string
}

type Tool struct {
	toolID			int
	toolName		string
	toolOrigName	string
}


/**
	Check if a text answer is "yes"
 */
func EvalYES(text string) bool {
	text = strings.Trim(text, " ")
	if strings.EqualFold(text, "YES") || strings.EqualFold(text, "ok") {
		return true
	}

	return false
}

func ExtractToolDescriptionAndTags(text string) (AlertRequestItems) {
	return AlertRequestItems{}
}
