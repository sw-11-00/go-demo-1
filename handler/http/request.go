package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string, result interface{}) error {
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
