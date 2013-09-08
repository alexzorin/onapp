package cmd

import (
	"flag"
	"fmt"
	"os"
)

type cli struct {
	*config
}

type cmdHandler interface {
	Run([]string, *cli) error
	Description() string
}

var cmdHandlers = map[string]cmdHandler{
	"config": configCmd{},
}

func (c *cli) parse(args []string) {
	if len(args) == 0 {
		errorln("\nNo command passed, usage is as follows:")
		printUsage()
		return
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		err := handler.Run(args[1:], c)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		errorf("\nERROR: %s is an unknown command\n", args[0])
		printUsage()
	}
}

func Start() {
	infoln("===> Starting OnApp CLI ...\n")
	conf, err := loadConfig()
	if err != nil {
		errorf("=== ERROR\n", err)
		os.Exit(1)
	}
	cli := cli{conf}
	cli.parse(cleanArgs(os.Args[1:]))
}

func printUsage() {
	infoln("\n===> USAGE")
	infoln("Available commands:")
	for k, v := range cmdHandlers {
		infof("  %s - %s\n", k, v.Description())
	}
	infoln("\nGeneral options:")
	infoToggle(true)
	flag.PrintDefaults()
	infoToggle(false)
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
