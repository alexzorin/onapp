package cmd

import (
	"flag"
	"fmt"
	"github.com/alexzorin/onapp/cmd/log"
	"os"
	"path/filepath"
)

type cli struct {
	*config
	caller string
}

type cmdHandler interface {
	Run([]string, *cli) error
	Description() string
	Help([]string)
}

var cmdHandlers = map[string]cmdHandler{
	"config": configCmd{},
	"help":   helpCmd{},
}

func (c *cli) parse(args []string) {
	if len(args) == 0 {
		log.Errorln("\nNo command passed")
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
	cli := cli{conf, filepath.Base(os.Args[0])}
	cli.parse(cleanArgs(os.Args[1:]))
}

func printUsage() {
	log.Infoln("\nAvailable commands\n")
	for k, v := range cmdHandlers {
		log.Infof("  %10s   %s\n", k, v.Description())
	}
	log.Infoln("\nGeneral options\n")
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
