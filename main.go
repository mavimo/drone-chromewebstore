package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Version set at compile-time
var (
	Version  string
	BuildNum string
)

func main() {
	app := cli.NewApp()
	app.Name = "Drone Chrome Webstore"
	app.Usage = "Deploying Chrome plugins on webstore with drone CI to an existing plugin"
	app.Copyright = "Copyright (c) 2018 Marco Vito Moscaritolo"
	app.Authors = []cli.Author{
		{
			Name:  "Marco Vito Moscaritolo",
			Email: "mavimo@gmail.com",
		},
	}
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env-file",
			Usage: "Env file to load (useful for debugging)",
		},
		cli.StringFlag{
			Name:   "application",
			Usage:  "Application ID",
			EnvVar: "PLUGIN_APPLICATION",
		},
		cli.StringFlag{
			Name:   "client-id",
			Usage:  "Client ID",
			EnvVar: "PLUGIN_CLIENT_ID",
		},
		cli.StringFlag{
			Name:   "client-secret",
			Usage:  "Client secret",
			EnvVar: "PLUGIN_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "refresh-token",
			Usage:  "Refresh token",
			EnvVar: "PLUGIN_REFRESH_TOKEN",
		},
		cli.StringFlag{
			Name:   "source",
			Usage:  "Application source folder",
			EnvVar: "PLUGIN_SOURCE",
		},
		cli.BoolTFlag{
			Name:   "upload",
			Usage:  "Upload application to webstore",
			EnvVar: "PLUGIN_UPLOAD",
		},
		cli.BoolTFlag{
			Name:   "publish",
			Usage:  "Pubplish application in webstore",
			EnvVar: "PLUGIN_PUBLISH",
		},
		cli.StringFlag{
			Name:   "publish-target",
			Usage:  "Publish target, should be default or trustedTesters",
			EnvVar: "PLUGIN_PUBLISH_TARGET",
			Value:  "default",
		},
	}

	app.Version = Version

	if BuildNum != "" {
		app.Version = app.Version + "+" + BuildNum
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Warningln(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		ApplicationID: c.String("application"),
		Authentication: Authentication{
			ClientID:     c.String("client-id"),
			ClientSecret: c.String("client-secret"),
			RefreshToken: c.String("refresh-token"),
		},
		Config: Config{
			Source:        c.String("source"),
			Upload:        c.BoolT("upload"),
			Publish:       c.BoolT("publish"),
			PublishTarget: c.String("publish-target"),
		},
	}

	return plugin.Exec()
}
