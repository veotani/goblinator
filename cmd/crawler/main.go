package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/veotani/goblinator/pkg/config"
)

type BlizzardTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Sub         string `json:"sub"`
}

func getConfig() *config.GoblinatorConfig {
	config, err := config.New()
	if err != nil {
		log.Fatalf("Couldn't set up configuration: %s", err)
	}
	return config
}

func authRequest(clientId string, clientSecret string) (*http.Request, error) {
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://eu.battle.net/oauth/token", body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func parseAuthResp(authResp *http.Response) (*BlizzardTokenResponse, error) {
	token := &BlizzardTokenResponse{}
	respBody, err := ioutil.ReadAll(authResp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respBody, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func auth(clientId string, clientSecret string) (*BlizzardTokenResponse, error) {
	req, err := authRequest(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Request wasn't successfull and returned code %d", resp.StatusCode)
	}

	token, err := parseAuthResp(resp)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func main() {
	config := getConfig()
	token, err := auth(config.BlizzardClientId, config.BlizzardClientSecret)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	log.Printf("Authentication success: %v", token)
}
