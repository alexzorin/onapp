package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/log"
	"os"
	"path/filepath"
)

type cli struct {
	*config
	caller    string
	apiClient *onapp.Client
}

type cmdHandler interface {
	Run([]string, *cli) error
	Description() string
}

type cmdHandlerHelp interface {
	Help([]string)
}

type cmdHandlerSubhandlers interface {
	Handlers() *map[string]cmdHandler
}

var cmdHandlers = map[string]cmdHandler{
	"config": configCmd{},
	"vm":     vmCmd{},
	"test":   testCmd{},
	"help":   helpCmd{},
}

func (c *cli) parse(args []string) {
	if len(args) == 0 {
		log.Errorln("No command passed")
		printUsage()
		return
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		err := handler.Run(args[1:], c)
		if err != nil {
			log.Errorln(err)
		}
	} else {
		log.Errorf("%s is an unknown command\n", args[0])
		printUsage()
	}
}

func Start() {
	conf, err := loadConfig()
	if err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
	cl, err := onapp.NewClient(conf.Server, conf.ApiUser, conf.ApiKey)
	if err != nil {
		log.Errorf(err.Error())
	}
	cli := cli{conf, filepath.Base(os.Args[0]), cl}
	cli.parse(cleanArgs(os.Args[1:]))
}

func (c *cli) subhandle(handler cmdHandlerSubhandlers, args []string) error {
	sub, ok := (*handler.Handlers())[args[0]]
	if ok {
		return sub.Run(args[1:], c)
	}
	return errors.New(fmt.Sprintf("Sub-command %s doesn't exist", args[0]))
}

func printUsage() {
	log.Infoln("Available commands\n")
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
