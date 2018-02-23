// Functions to simplify HTTP requests

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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
