package cmd

import (
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/log"
	cgcli "github.com/codegangsta/cli"
)

func NewTestCmd(c *cli) cgcli.Command {
	return cgcli.Command{
		Name:  "test",
		Usage: "Try to connect to the dashboard server",
		Flags: []cgcli.Flag{
			cgcli.StringFlag{
				Name:  "host",
				Value: c.config.Server,
				Usage: "API server URL",
			},
			cgcli.StringFlag{
				Name:  "apiUser",
				Value: c.config.ApiUser,
				Usage: "API User (usually email address)",
			},
			cgcli.StringFlag{
				Name:  "apiKey",
				Value: c.config.ApiKey,
				Usage: "API key",
			},
		},
		Action: func(ctx *cgcli.Context) {
			client, err := onapp.NewClient(ctx.String("host"), ctx.String("apiUser"), ctx.String("apiKey"))
			if err != nil {
				log.Errorf("Failed to build client: %s\n", err.Error())
				return
			}
			p, err := client.GetProfile()
			if err != nil {
				log.Errorf("Failed to get profile: %s\n", err.Error())
				return
			}
			log.Successf("Successfully connected! Your profile name is %s %s.\n", p.FirstName, p.LastName)
		},
	}
}
