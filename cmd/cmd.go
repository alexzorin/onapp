package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/log"
	cgcli "github.com/codegangsta/cli"
)

type cli struct {
	*config
	caller    string
	apiClient *onapp.Client
	cache     Cache
	app       *cgcli.App
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
	"help":   helpCmd{},
}

func (c *cli) setup() {
	c.app = cgcli.NewApp()
	c.app.Name = "onapp"
	c.app.Usage = "Interact with the OnApp API"
	c.app.Commands = []cgcli.Command{
		NewTestCmd(c),
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
	cli := cli{config: conf, caller: filepath.Base(os.Args[0]), apiClient: cl, cache: &fileBackedCache{}}
	cli.setup()
	cli.app.Run(os.Args)
}

func (c *cli) subhandle(handler cmdHandlerSubhandlers, args []string) error {
	sub, ok := (*handler.Handlers())[args[0]]
	if ok {
		return sub.Run(args[1:], c)
	}
	return errors.New(fmt.Sprintf("Sub-command %s doesn't exist", args[0]))
}
