package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/veotani/goblinator/pkg/config"
)

func getConfig() *config.GoblinatorConfig {
	config, err := config.New()
	if err != nil {
		log.Fatalf("Couldn't set up configuration: %s", err)
	}
	return config
}

func authenticate(clientId string, clientSecret string) *http.Response {
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://eu.battle.net/oauth/token", body)
	if err != nil {
		log.Fatalf("Failed to create request: %s", err)
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Request wasn't successfull and returned code %d", resp.StatusCode)
	}

	return resp
}

func main() {
	config := getConfig()
	resp := authenticate(config.BlizzardClientId, config.BlizzardClientSecret)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't read response: %s", err)
	}
	log.Println(string(respBody))
}
