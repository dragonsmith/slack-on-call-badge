// Sets "On call" badge to a Slack user it one is on call in OpsGenie rotation.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type slackStatus struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type genieRorationJSON struct {
	Participants []struct {
		Name string `json:"name"`
	} `json:"participants"`
}

type slackProfileJSON struct {
	OK      bool        `json:"ok,omitempty"`
	Error   string      `json:"error,omitempty"`
	User    string      `json:"user"`
	Profile slackStatus `json:"profile"`
}

type adminAccount struct {
	slackID     string
	genieOnCall bool
	slackOnCall bool
}

const (
	opsGenieAPIURL = "https://api.opsgenie.com/v1.1/json/schedule/whoIsOnCall?apiKey=%s&name=%s"
	slackAPIURL    = "https://slack.com/api/users.profile.get?token=%s&user=%s"
)

var (
	slackToken   = os.Getenv("SLACK_TOKEN")
	genieToken   = os.Getenv("OPSGENIE_TOKEN")
	rotationName = os.Getenv("OPSGENIE_ROTATION")
	adminsEnv    = os.Getenv("ADMINS")

	admins = parseUsers(adminsEnv)
)

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
}

func parseUsers(users string) map[string]adminAccount {
	parsedUsers := make(map[string]adminAccount)

	for _, user := range strings.Split(users, ",") {

		splitUser := strings.Split(user, ":")

		if len(splitUser) != 2 {
			err := "\"ADMINS\" variable contain malformed string.\n"
			err += "Please format is like that:"
			err += "\tUser1_OpsGenie_email:User1_Slack_id,User2_OpsGenie_email:User2_Slack_id,...\n"
			log.Fatalln(err)
		}

		parsedUsers[splitUser[0]] = adminAccount{slackID: splitUser[1]}
	}

	return parsedUsers
}

func httpGet(url string) ([]byte, error) {
	var response []byte

	resp, err := http.Get(url)
	if err != nil {
		return response, err
	}

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = resp.Body.Close()
	if err != nil {
		return response, err
	}

	return response, err
}

func httpPostJSON(url string, data interface{}, headers map[string]string) error {
	var slackClient = http.Client{}

	newProfileJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, reqErr := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(newProfileJSON))
	if reqErr != nil {
		return reqErr
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	_, postErr := slackClient.Do(req)
	return postErr
}

func whoIsOnCallOpsGenie(token string, schedule string, admins map[string]adminAccount) {
	var parsedJSON genieRorationJSON

	apiURL := fmt.Sprintf(opsGenieAPIURL, url.QueryEscape(token), url.QueryEscape(schedule))

	unparsedJSON, err := httpGet(apiURL)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(unparsedJSON, &parsedJSON)
	if err != nil {
		log.Fatalln(err)
	}

	for _, arg := range parsedJSON.Participants {
		email := arg.Name
		if data, ok := admins[email]; ok {
			data.genieOnCall = true
			admins[email] = data
		}
	}
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

		if adminStatus.StatusEmoji == ":on_call:" {
			data.slackOnCall = true
			admins[email] = data
		}
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

	postErr := httpPostJSON(apiURL, newProfile, headers)
	if postErr != nil {
		log.Fatalln(postErr)
	}
}

func main() {
	checkConfig()

	log.Print("Daemon started!\n\n")

	log.Println("Managing badges for:")
	for email := range admins {
		log.Println("* " + email)
	}
	log.Println("")

	for {
		go func() {
			// Fill OpsGenie on duty flag inside "admins" map
			whoIsOnCallOpsGenie(genieToken, rotationName, admins)

			// Fill Slack on duty flag inside "admins" map
			whoIsOnCallSlack(slackToken, admins)

			for email, data := range admins {
				if data.genieOnCall && !data.slackOnCall {
					log.Println("Setting badge for user:", email)
					go setSlackStatus(slackToken, data.slackID, ":on_call:", "on call")
				}

				if !data.genieOnCall && data.slackOnCall {
					log.Println("Unsetting badge for user:", email)
					go setSlackStatus(slackToken, data.slackID, "", "")
				}
			}
		}()

		time.Sleep(10 * time.Second)
	}
}
