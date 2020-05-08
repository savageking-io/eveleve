package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var AppVersion = "Undefined"
var BuildID = "Undefined"

func main() {
	app := cli.NewApp()
	app.Name = "eveleve"
	app.Version = AppVersion
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Mike Savage King",
			Email: "mike@savageking.io",
		},
	}
	app.Description = "Manage your Indie Gamedev Community"
	app.Copyright = "Copyright 2020 Mike Savage King"

	app.Commands = []*cli.Command{
		{
			Name:  "master",
			Usage: "Run EvelEve in Master mode",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) error {
				log.SetLevel(log.TraceLevel)
				var m Master
				if err := m.Init(); err != nil {
					log.Errorf("Failed to initialize master server: %s", err.Error())
					return err
				}
				return m.Run()
			},
		},
	}
	app.Run(os.Args)
}
