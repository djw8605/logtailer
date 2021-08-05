// Package device facilitates performing OAuth Device Authorization Flow for client applications
// such as CLIs that can not receive redirects from a web site.
//
// First, RequestCode should be used to obtain a CodeResponse.
//
// Next, the user will need to navigate to VerificationURI in their web browser on any device and fill
// in the UserCode.
//
// While the user is completing the web flow, the application should invoke PollToken, which blocks
// the goroutine until the user has authorized the app on the server.
//
// https://docs.github.com/en/free-pro-team@latest/developers/apps/authorizing-oauth-apps#device-flow
package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnsupported is thrown when the server does not implement Device flow.
	ErrUnsupported = errors.New("device flow not supported")
	// ErrTimeout is thrown when polling the server for the granted token has timed out.
	ErrTimeout = errors.New("authentication timed out")
)

type CodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	timeNow         func() time.Time
	timeSleep       func(time.Duration)
}

type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// RequestCode initiates the authorization flow by requesting a code from uri.
func RequestCode(uri string, clientID string, scopes []string) (*CodeResponse, error) {

	resp, err := http.PostForm(uri, url.Values{
		"client_id": {clientID},
		"scope":     {strings.Join(scopes, " ")},
	})
	if err != nil {
		return nil, err
	}

	// Parse the json returned
	var bb []byte
	bb, err = ioutil.ReadAll(resp.Body)
	responseDevice := CodeResponse{}
	log.Errorln("Received:", string(bb))
	if err := json.Unmarshal(bb, &responseDevice); err != nil {
		log.Fatalln("Unable to unmarshall device json:", err)
	}

	return &responseDevice, nil
}

var grantType string = "urn:ietf:params:oauth:grant-type:device_code"

// PollToken polls the server at pollURL until an access token is granted or denied.
func PollToken(pollURL string, clientID string, code *CodeResponse) (*ResponseToken, error) {
	timeNow := code.timeNow
	if timeNow == nil {
		timeNow = time.Now
	}
	timeSleep := code.timeSleep
	if timeSleep == nil {
		timeSleep = time.Sleep
	}

	checkInterval := time.Duration(code.Interval) * time.Second
	expiresAt := timeNow().Add(time.Duration(code.ExpiresIn) * time.Second)

	for {
		timeSleep(checkInterval)

		resp, err := http.PostForm(pollURL, url.Values{
			"client_id":   {clientID},
			"device_code": {code.DeviceCode},
			"grant_type":  {grantType},
		})
		if err != nil {
			log.Errorln("Post form failed:", err)
			return nil, err
		}

		if resp.StatusCode != 200 {
			log.Debugln("Status Code:", resp.StatusCode)
			var bb []byte
			bb, err = ioutil.ReadAll(resp.Body)
			log.Debugln(string(bb))
		} else {
			var bb []byte
			bb, err = ioutil.ReadAll(resp.Body)
			response := ResponseToken{}
			if err := json.Unmarshal(bb, &response); err != nil {
				log.Fatalln("Unable to unmarshall device json:", err)
			}
			return &response, nil

		}

		if timeNow().After(expiresAt) {
			return nil, ErrTimeout
		}
	}
}
