package cmd

import (
	"github.com/alexzorin/onapp/cmd/log"
)

const (
	vmCmdDescription     = "Manage virtual machines"
	vmCmdHelp            = "See subcommands for help on managing virtual machines."
	vmCmdListDescription = "List virtual machines under your account"
	vmCmdListHelp        = "\nUsage: `onapp vm list [filter]`\n" +
		"Optionally filter by virtual machine label.\n"
)

// Base command

type vmCmd struct{}

var vmCmdHandlers = map[string]cmdHandler{
	"list": vmCmdList{},
}

func (c vmCmd) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		log.Infoln("This command does nothing when invoked on its own.")
		cmdHandlers["help"].Run([]string{"vm"}, ctx)
		return nil
	} else {
		return ctx.subhandle(c, args)
	}
}

func (c vmCmd) Description() string {
	return vmCmdDescription
}

func (c vmCmd) Help(args []string) {
	log.Infoln(vmCmdHelp)
}

func (c vmCmd) Handlers() *map[string]cmdHandler {
	return &vmCmdHandlers
}

// List command
type vmCmdList struct{}

func (c vmCmdList) Run(args []string, ctx *cli) error {
	return nil
}

func (c vmCmdList) Description() string {
	return vmCmdListDescription
}

func (c vmCmdList) Help(args []string) {
	log.Infoln(vmCmdListHelp)
}
