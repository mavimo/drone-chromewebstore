package main

import (
	"fmt"
)

// Plugin to deploy application in chrome webstore
type Plugin struct {
	ApplicationID  string
	Config         Config
	Authentication Authentication
}

// Authentication contains settings required to authenticate API
type Authentication struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

// Config indication operation to do in plugin
type Config struct {
	Source        string
	Upload        bool
	Publish       bool
	PublishTarget string
}

// Exec operation for this plugin
func (p Plugin) Exec() error {
	client := NewChromeWebstoreClient(p.ApplicationID, p.Authentication)

	if p.Config.Upload {
		buf, err := GenerateZipContent(p.Config.Source)
		if err != nil {
			return fmt.Errorf("unable to generate zip content: %v", err)
		}

		if err := client.UploadNewVersion(buf); err != nil {
			return fmt.Errorf("unable to upload a new versio: %v", err)
		}
	}

	if p.Config.Publish {
		if err := client.PublishVersion(p.Config.PublishTarget); err != nil {
			return fmt.Errorf("unable to publish a new version: %v", err)
		}
	}

	return nil
}
