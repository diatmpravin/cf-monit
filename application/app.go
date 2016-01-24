package app

import (
	"cf-appmonit/model"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func getToken(domain, username, password string) (response model.AuthenticationResponse) {
	uaaUrl := fmt.Sprintf("http://uaa.%s/oauth/token", domain)

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)
	data.Add("grant_type", "password")
	data.Add("scope", "")

	request, err := http.NewRequest("POST", uaaUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")))

	request.Header.Set("accept", "application/json")

	client := newClient()
	rawResponse, err := client.Do(request)
	if err != nil {
		log.Println("Error performing request: %s", err.Error())
	}

	if rawResponse.StatusCode > 299 {
		log.Println("Server error, status code: %d", rawResponse.StatusCode)
	}

	jsonBytes, err := ioutil.ReadAll(rawResponse.Body)
	rawResponse.Body.Close()

	if err != nil {
		log.Println("Could not read response body: %s", err.Error())
	}

	err = json.Unmarshal(jsonBytes, &response)

	if err != nil {
		log.Println("Invalid JSON response from server: %s", err.Error())
	}

	return response
}

func getappGUID(accessToken, domain, appName string) (appGUID string) {
	uaaUrl := fmt.Sprintf("http://api.%s/v2/apps", domain)
	bearer := fmt.Sprintf("bearer %s", accessToken)

	request, err := http.NewRequest("GET", uaaUrl, nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", bearer)

	client := newClient()
	rawResponse, err := client.Do(request)
	if err != nil {
		log.Printf("Error performing request: %s", err.Error())
	}

	if rawResponse.StatusCode > 299 {
		log.Println("Server error, status code: %d", rawResponse.StatusCode)
	}

	jsonBytes, err := ioutil.ReadAll(rawResponse.Body)
	rawResponse.Body.Close()

	if err != nil {
		log.Println("Could not read response body: %s", err.Error())
	}

	response := model.ApiResponse{}
	err = json.Unmarshal(jsonBytes, &response)

	if err != nil {
		log.Println("Invalid JSON response from server: %s", err.Error())
	}

	for _, v := range response.Resources {
		if v.Entity.Name == appName {
			appGUID = v.Metadata.Guid
		}
	}

	return appGUID
}

func displayStats(appGUID, domain, appName, token string) {
	for {
		uaaUrl := fmt.Sprintf("http://api.%s/v2/apps/%s/stats", domain, appGUID)
		bearer := fmt.Sprintf("bearer %s", token)

		request, err := http.NewRequest("GET", uaaUrl, nil)
		if err != nil {
			return
		}
		request.Header.Set("Authorization", bearer)

		client := newClient()
		rawResponse, err := client.Do(request)
		if err != nil {
			log.Printf("Error performing request: %s", err.Error())
		}

		if rawResponse.StatusCode > 299 {
			log.Println("Server error, status code: %d", rawResponse.StatusCode)
		}

		jsonBytes, err := ioutil.ReadAll(rawResponse.Body)
		rawResponse.Body.Close()

		if err != nil {
			log.Println("Could not read response body: %s", err.Error())
		}

		response := model.AppStat{}
		err = json.Unmarshal(jsonBytes, &response.Data)

		if err != nil {
			log.Println("Invalid JSON response from server: %s", err.Error())
		}

		for _, v := range response.Data {
			log.Println("------State----->", v.State)
			log.Println("------Name----->", v.Stats.Name)
			log.Println("------Up Time----->", v.Stats.Usage.Time)
			log.Println("------Memory----->", v.Stats.Usage.Memory)
			log.Println("------CPU----->", v.Stats.Usage.CPU)
			log.Println("------Disk----->", v.Stats.Usage.Disk)
			log.Println("\n\n")
		}
		time.Sleep(time.Second * 2)
	}

}

func Monitor(domain, username, password, appName string) {
	tokenStr := getToken(domain, username, password)
	appGUID := getappGUID(tokenStr.AccessToken, domain, appName)
	displayStats(appGUID, domain, appName, tokenStr.AccessToken)
}
