package cmd

import (
	"flag"
	"fmt"
	"github.com/alexzorin/onapp/cmd/log"
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
		log.Errorln("\nNo command passed, usage is as follows:")
		printUsage()
		return
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		err := handler.Run(args[1:], c)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		log.Errorf("\nERROR: %s is an unknown command\n", args[0])
		printUsage()
	}
}

func Start() {
	conf, err := loadConfig()
	if err != nil {
		log.Errorf("=== ERROR\n", err)
		os.Exit(1)
	}
	cli := cli{conf}
	cli.parse(cleanArgs(os.Args[1:]))
}

func printUsage() {
	log.Infoln("\n===> USAGE")
	log.Infoln("Available commands:")
	for k, v := range cmdHandlers {
		log.Infof("  %s - %s\n", k, v.Description())
	}
	log.Infoln("\nGeneral options:")
	log.InfoToggle(true)
	flag.PrintDefaults()
	log.InfoToggle(false)
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
