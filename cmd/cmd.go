package cmd

import (
	"flag"
	"fmt"
	"os"
)

type CLI struct {
	*config
}

type cmdHandler interface {
	Run([]string, *CLI) error
	Description() string
}

var cmdHandlers = map[string]cmdHandler{
	"config": configCmd{},
}

func (c *CLI) parse(args []string) {
	if len(args) == 0 {
		Infoln("\nNo command passed, usage is as follows:")
		printUsage()
		return
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		err := handler.Run(args[1:], c)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		Errorf("\nERROR: %s is an unknown command\n", args[0])
		printUsage()
	}
}

func Start() {
	Infoln("===> Starting OnApp CLI ...\n")
	conf, err := LoadConfig()
	if err != nil {
		Errorf("=== ERROR\n", err)
		os.Exit(1)
	}
	cli := CLI{conf}

	cli.parse(cleanArgs(os.Args[1:]))
}

func printUsage() {
	Infoln("\n===> USAGE")
	Infoln("Available commands:")
	for k, v := range cmdHandlers {
		Infof("  %s - %s\n", k, v.Description())
	}
	Infoln("\nGeneral options:")
	InfoToggle(true)
	flag.PrintDefaults()
	InfoToggle(false)
}

func cleanArgs(args []string) []string {
	out := make([]string, 0)
	for _, v := range args {
		if v[0] != '-' {
			out = append(out, v)
		}
	}
	return out
}
