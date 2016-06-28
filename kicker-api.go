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
			Value:   "",
			Usage:   "Datasource connection string",
			EnvVars: []string{"KICKER_DSN"},
		},
		&cli.StringFlag{
			Name:    "dialect",
			Value:   "mysql",
			Usage:   "Datasource dialect",
			EnvVars: []string{"KICKER_DB_DIALECT"},
		},
	}

	app.Action = runServer

	app.Run(os.Args)
}
