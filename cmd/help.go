package cmd

import (
	"github.com/alexzorin/onapp/cmd/log"
)

const (
	helpCmdDescription = "Help text for subcommands"
	helpCmdHelp        = "To get help with a command, use `help [command]`"
)

type helpCmd struct {
}

func (c helpCmd) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		handler.Help(args[1:])
	} else {
		log.Errorf("Command '%s' not found\n", args[0])
	}
	return nil
}

func (c helpCmd) Description() string {
	return helpCmdDescription
}

func (c helpCmd) Help(args []string) {
	log.Infoln(helpCmdHelp)
	printUsage()
}
