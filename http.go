// Functions to simplify HTTP requests

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet(url string, headers map[string]string) ([]byte, error) {
	var (
		body []byte
		resp *http.Response
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return body, err
	}

	for name, value := range headers {
		req.Header.Add(name, value)
	}

	resp, err = client.Do(req)
	if err != nil {
		return body, err
	}

	if resp.StatusCode != 200 {
		errormsg := "Error during HTTP GET to OpsGenie API: %d %s"
		return body, errors.New(fmt.Sprintf(errormsg, resp.StatusCode, http.StatusText(resp.StatusCode)))
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	err = resp.Body.Close()

	return body, err
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
