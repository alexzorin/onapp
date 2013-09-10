package cmd

import (
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/cmd/log"
)

const (
	testCmdDescription = "Try to connect to the dashboard server"
	testCmdHelp        = "This command will try to connect to the dashboard server and access your user profile"
)

type testCmd struct {
}

func (c testCmd) Run(args []string, ctx *cli) error {
	client, err := onapp.NewClient(ctx.config.Server, ctx.config.ApiUser, ctx.config.ApiKey)
	if err != nil {
		return err
	}
	p, err := client.GetProfile()
	if err != nil {
		return err
	}
	log.Infof("Successfully connected! Your profile name is %s %s.\n", p.FirstName, p.LastName)
	return nil
}

func (c testCmd) Description() string {
	return testCmdDescription
}

func (c testCmd) Help(args []string) {
	log.Infoln(testCmdHelp)
}
