package cmd

import (
	"errors"
	"fmt"
	"github.com/alexzorin/onapp/log"
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
		// Print help text if available
		helping, ok := handler.(cmdHandlerHelp)
		if ok {
			helping.Help(args[1:])
		}
		// List subcommands if available
		subhandled, ok := handler.(cmdHandlerSubhandlers)
		if ok {
			handlers := subhandled.Handlers()
			// If the user passed the subhandler to the help command
			// and the subhandler satisfies the help interface
			// call help on the subhandler too.
			if len(args) > 1 {
				subhandler, ok := (*handlers)[args[1]]
				if ok {
					helping, ok = subhandler.(cmdHandlerHelp)
					if ok {
						helping.Help(args[2:])
					}
				}
			}
			// and list subhandlers at the end
			log.Infof("`%s' has a number of sub-commands:\n", args[0])
			for k, v := range *handlers {
				log.Infof("  %10s   %s\n", k, v.Description())
			}
		}
	} else {
		return errors.New(fmt.Sprintf("Command '%s' not found", args[0]))
	}
	return nil
}

func (c helpCmd) Description() string {
	return helpCmdDescription
}

func (c helpCmd) Help(args []string) {
	log.Infoln(helpCmdHelp, "\n")
	printUsage()
}
