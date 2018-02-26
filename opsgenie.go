// Functions to get information from OpsGenie

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const (
	opsGenieAPIURL = "https://api.opsgenie.com/v1.1/json/schedule/whoIsOnCall?apiKey=%s&name=%s"
)

type genieRorationJSON struct {
	Participants []struct {
		Name string `json:"name"`
	} `json:"participants"`
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

	if *debug {
		log.Println("DEBUG: Raw data from OpsGenie")
		log.Println("DEBUG:", parsedJSON.Participants)
	}

	for emailFromConfig, dataFromConfig := range admins {
		oncall := false

		for _, adminFromGenie := range parsedJSON.Participants {
			if emailFromConfig == adminFromGenie.Name {
				oncall = true
			}
		}

		dataFromConfig.genieOnCall = oncall
		admins[emailFromConfig] = dataFromConfig
	}
}
