// Functions to get and set status in Slack

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

const (
	slackAPIURL = "https://slack.com/api/users.profile.get?token=%s&user=%s"
)

type slackStatus struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type slackProfileJSON struct {
	OK      bool        `json:"ok,omitempty"`
	Error   string      `json:"error,omitempty"`
	User    string      `json:"user"`
	Profile slackStatus `json:"profile"`
}

func getSlackUserStatus(token string, user string) (slackStatus, error) {
	var parsedJSON slackProfileJSON
	apiURL := fmt.Sprintf(slackAPIURL, url.QueryEscape(token), url.QueryEscape(user))

	unparsedJSON, err := httpGet(apiURL)
	if err != nil {
		return slackStatus{"err", "err"}, err
	}

	err = json.Unmarshal(unparsedJSON, &parsedJSON)
	if err != nil {
		return slackStatus{"err", "err"}, err
	}

	if status := parsedJSON.OK; status {
		return slackStatus{parsedJSON.Profile.StatusText, parsedJSON.Profile.StatusEmoji}, nil
	}

	errText := fmt.Sprintf("Slack returned: %s; SlackID: %s;", parsedJSON.Error, user)
	return slackStatus{"err", "err"}, errors.New(errText)
}

func whoIsOnCallSlack(token string, admins map[string]adminAccount) {
	for email, data := range admins {
		adminStatus, err := getSlackUserStatus(token, data.slackID)
		if err != nil {
			log.Fatalln(err)
		}

		data.slackOnCall = (adminStatus.StatusEmoji == onCallIcon)
		admins[email] = data
	}
}

func setSlackStatus(token string, user string, badge string, text string) {

	var (
		apiURL      = "https://slack.com/api/users.profile.set"
		newProfile  = slackProfileJSON{User: user, Profile: slackStatus{text, badge}}
		bearerToken = fmt.Sprintf("Bearer %s", token)
	)

	headers := make(map[string]string)

	headers["Content-type"] = "application/json;charset=utf-8"
	headers["Authorization"] = bearerToken

	if *dryRun {
		log.Println("Dummy http request to set Slack Status:", user, badge, text)
	} else {
		postErr := httpPostJSON(apiURL, newProfile, headers)
		if postErr != nil {
			log.Fatalln(postErr)
		}
	}
}
