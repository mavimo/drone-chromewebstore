package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

// ChromeWebstoreClient create an http client to interact with Crome Webstore API
type ChromeWebstoreClient struct {
	*http.Client
	ApplicationID string
}

// NewChromeWebstoreClient generate a new client to interact with Chrome Webstore API
func NewChromeWebstoreClient(applicationID string, auth Authentication) ChromeWebstoreClient {
	ctx := context.Background()
	cfg := oauth2.Config{
		ClientID:     auth.ClientID,
		ClientSecret: auth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Scopes: []string{
			"https://www.googleapis.com/auth/chromewebstore",
		},
	}

	ts := cfg.TokenSource(ctx, &oauth2.Token{
		RefreshToken: auth.RefreshToken,
	})

	tkn, err := ts.Token()
	if err != nil {
		fmt.Println("unable to refresh token")
	}

	return ChromeWebstoreClient{cfg.Client(ctx, tkn), applicationID}
}

// UploadNewVersion send a new version of application to Chrome Webstore
func (client ChromeWebstoreClient) UploadNewVersion(buf *bytes.Buffer) error {
	// Try to upload zip file
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://www.googleapis.com/upload/chromewebstore/v1.1/items/%s", client.ApplicationID), buf)
	if err != nil {
		return fmt.Errorf("unable to create upload request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-goog-api-version", "2")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to upload zip file: %v", err)
	}
	defer res.Body.Close()

	message, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to get response when upload application %s: %v", client.ApplicationID, err)
	}

	log.Printf(string(message))

	return nil
}

// GetInfo get information on an application froom Chrome Webstore
func (client ChromeWebstoreClient) GetInfo() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/chromewebstore/v1.1/items/%s?projection=DRAFT", client.ApplicationID), nil)
	if err != nil {
		return fmt.Errorf("unable to create info request: %v", err)
	}
	req.Header.Set("x-goog-api-version", "2")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to fetch infos for application %s: %v", client.ApplicationID, err)
	}
	defer res.Body.Close()

	message, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to get response when get info for application %s: %v", client.ApplicationID, err)
	}

	log.Printf(string(message))

	return nil
}

// PublishVersion get information on an application froom Chrome Webstore
func (client ChromeWebstoreClient) PublishVersion(target string) error {
	if target != "default" && target != "trustedTesters" {
		return fmt.Errorf("unable to publish application %s", client.ApplicationID)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://www.googleapis.com/chromewebstore/v1.1/items/%s/publish", client.ApplicationID), nil)
	if err != nil {
		return fmt.Errorf("unable to create publish request: %v", err)
	}
	req.Header.Set("x-goog-api-version", "2")
	req.URL.Query().Add("publishTarget", target)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to publish application %s: %v", client.ApplicationID, err)
	}
	defer res.Body.Close()

	message, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to get response when publish application %s: %v", client.ApplicationID, err)
	}

	log.Printf(string(message))

	return nil
}
