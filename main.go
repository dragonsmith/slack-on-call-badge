// Sets "On call" badge to a Slack user it one is on call in OpsGenie rotation.
package main

import (
	"log"
	"os"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	slackToken   = os.Getenv("SLACK_TOKEN")
	genieToken   = os.Getenv("OPSGENIE_TOKEN")
	rotationName = os.Getenv("OPSGENIE_ROTATION")
	adminsEnv    = os.Getenv("ADMINS")

	onCallIcon = os.Getenv("ON_CALL_ICON")
	onCallText = os.Getenv("ON_CALL_TEXT")

	admins = parseUsers(adminsEnv)

	runOnce = kingpin.Flag("once", "Run once instead staying in foreground with periodic checks").Bool()
)

func checkAndUpdate() {
	// Fill OpsGenie on duty flag inside "admins" map
	whoIsOnCallOpsGenie(genieToken, rotationName, admins)

	// Fill Slack on duty flag inside "admins" map
	whoIsOnCallSlack(slackToken, admins)

	for email, data := range admins {
		if data.genieOnCall && !data.slackOnCall {
			log.Println("Setting badge for user:", email)
			setSlackStatus(slackToken, data.slackID, onCallIcon, onCallText)
		}

		if !data.genieOnCall && data.slackOnCall {
			log.Println("Unsetting badge for user:", email)
			setSlackStatus(slackToken, data.slackID, "", "")
		}
	}
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	checkConfig()

	log.Println("Managing badges for:")
	for email := range admins {
		log.Println("* " + email)
	}
	log.Println("")

	if *runOnce {
		checkAndUpdate()
	} else {
		for {
			go checkAndUpdate()
			time.Sleep(60 * time.Second)
		}
	}
}
