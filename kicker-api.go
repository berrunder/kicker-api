package main

import (
	"gopkg.in/urfave/cli.v2"

	"os"
)

func main() {
	app := cli.NewApp()

	app.Name = "Kicker Api"
	app.Usage = "Kicker RESTFUL API"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   "4000",
			Usage:   "Exposed port",
		},
		&cli.StringFlag{
			Name:    "datasource",
			Aliases: []string{"d"},
			Value:   "kicker:starcode@tcp(localhost:3506)/kicker?parseTime=true",
			Usage:   "Datasource connection string (MySQL)",
			EnvVars: []string{"KICKER_DSN"},
		},
	}

	app.Action = runServer

	app.Run(os.Args)
}
