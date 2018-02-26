// Functions that handle configuration.

package main

import (
	"log"
	"strings"
)

type adminAccount struct {
	slackID     string
	genieOnCall bool
	slackOnCall bool
}

func checkConfig() {
	if slackToken == "" {
		log.Fatalln("\"SLACK_TOKEN\" variable should be defined!")
	}

	if genieToken == "" {
		log.Fatalln("\"OPSGENIE_TOKEN\" variable should be defined!")
	}

	if rotationName == "" {
		log.Fatalln("\"OPSGENIE_ROTATION\" variable should be defined!")
	}

	if adminsEnv == "" {
		log.Fatalln("\"ADMINS\" variable should be defined!")
	}

	// Set default values for on call text status and icon name
	if onCallIcon == "" {
		onCallIcon = ":on_call:"
	}

	if onCallText == "" {
		onCallText = "on call"
	}

	if *debug {
		log.Println("DEBUG: SLACK_TOKEN", slackToken)
		log.Println("DEBUG: OPSGENIE_TOKEN", genieToken)
		log.Println("DEBUG: OPSGENIE_ROTATION", rotationName)
		log.Println("DEBUG: ADMINS", adminsEnv)
		log.Println("DEBUG: ON_CALL_ICON", onCallIcon)
		log.Println("DEBUG: ON_CALL_TEXT", onCallText)
	}
}

func parseUsers(users string) map[string]adminAccount {
	parsedUsers := make(map[string]adminAccount)

	for _, user := range strings.Split(users, ",") {

		splitUser := strings.Split(user, ":")

		if len(splitUser) != 2 {
			log.Println("\"ADMINS\" variable contain malformed string.")
			log.Println("Please format is like that:")
			log.Fatalln("\tUser1_OpsGenie_email:User1_Slack_id,User2_OpsGenie_email:User2_Slack_id,...")
		}

		parsedUsers[splitUser[0]] = adminAccount{slackID: splitUser[1]}
	}

	return parsedUsers
}
