package main

import (
	"encoding/json"
	"fmt"
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

func fetchAuctionData(token string, realmId int) (*http.Response, error) {
	baseUrl := "https://eu.api.blizzard.com"

	params := url.Values{}
	// params.Add("region", "eu")
	params.Add("namespace", "dynamic-eu")
	params.Add("locale", "en_US")
	params.Add("access_token", token)
	log.Printf("Parameters string: %s", params.Encode())
	// body := strings.NewReader(params.Encode())

	endpoint := fmt.Sprintf(
		"/data/wow/connected-realm/%d/auctions",
		realmId,
	)
	log.Printf("Sending request to %s\n", baseUrl+endpoint)
	req, err := http.NewRequest(
		"GET",
		baseUrl+endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Send request to %s\n", resp.Request.URL)
	log.Printf("Status code is %d\n", resp.StatusCode)
	if resp.StatusCode != 200 {
		message, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Bad status without message: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("Bad status: %d. Message: %s", resp.StatusCode, string(message))
	}
	return resp, nil
}

func main() {
	config := getConfig()
	token, err := auth(config.BlizzardClientId, config.BlizzardClientSecret)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	log.Printf("Authentication success: %v", token)

	resp, err := fetchAuctionData(token.AccessToken, 509)
	if err != nil {
		log.Fatalln(err)
	}
	items, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(items))
}
