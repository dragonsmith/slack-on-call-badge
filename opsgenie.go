// Functions to get information from OpsGenie

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

const (
	opsGenieAPIURL = "https://api.opsgenie.com/v2/schedules/%s/on-calls?scheduleIdentifierType=name&flat=true"
)

type genieRorationJSON struct {
	Data struct {
		OnCallRecipients []string `json:"onCallRecipients"`
	} `json:"data"`
}

func whoIsOnCallOpsGenie(token string, schedule string, admins map[string]adminAccount) {
	var parsedJSON genieRorationJSON

	apiURL := fmt.Sprintf(opsGenieAPIURL, url.QueryEscape(schedule))

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("GenieKey %s", url.QueryEscape(token))

	unparsedJSON, err := httpGet(apiURL, headers)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(unparsedJSON, &parsedJSON)
	if err != nil {
		log.Fatalln(err)
	}

	if *debug {
		log.Println("DEBUG: Raw data from OpsGenie")
		log.Println("DEBUG:", parsedJSON.Data.OnCallRecipients)
	}

	for emailFromConfig, dataFromConfig := range admins {
		oncall := false

		for _, adminFromGenie := range parsedJSON.Data.OnCallRecipients {
			if emailFromConfig == adminFromGenie {
				oncall = true
			}
		}

		dataFromConfig.genieOnCall = oncall
		admins[emailFromConfig] = dataFromConfig
	}
}
