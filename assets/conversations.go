package assets

// -----------------------------------------------------------------------------------------
const IM_internal_error				= 0
const IM_introducing_myself			= 5
const IM_init_conversation 			= 6
const IM_accept_ticket_creation 	= 10
const IM_info_creating_issue		= 12
const IM_incident_created			= 15
const IM_what_can_i_do_for_you		= 80
const IM_error_cant_understand		= 99

// -----------------------------------------------------------------------------------------

const IM_event_join_channel =
	`Welcome to the channel. I'm iBot and you can contact me every time you detect an incident.`

// -----------------------------------------------------------------------------------------

const IM_asked_for_instructions =
	`These are the commands you can request to me.`

// -----------------------------------------------------------------------------------------

func GetSentence() func(int) string {
	// innerMap is captured in the closure returned below
	innerMap := map[int]string{
		IM_internal_error 			: 	`Oops, internal error. Could you give feedback to PES team so that they can fix this for the future?'`,
		IM_introducing_myself		: 	`Hi, my name is iBot and I'm here to create   ...alerts`,
		IM_init_conversation 		: 	`Do you want me to create jira incident on your behalf? [yes|no]`,
		IM_accept_ticket_creation 	:  	`Please type the incident you want me to create, here are some examples.
_Jira is down_
_Confluence is running slow_
_Jenkins. Some jobs are failing due to connection errors_
(The first word should be the name of the tool you want to report the problem from)`,
		IM_info_creating_issue		:   `I'm connecting with Jira server to create your issue. This can take up to a minute...`,
		IM_incident_created 		:	`Done. The alert has been created. Ticket: OCI-XXXX`,
		IM_what_can_i_do_for_you 	: 	`What can I do for you then?`,
		IM_error_cant_understand 	: 	`Sorry, I can't understand your answer.`,
	}

	return func(key int) string {
		return innerMap[key]
	}
}
