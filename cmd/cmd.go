package cmd

import (
	"fmt"
	"os"
)

type CLI struct {
	Handlers map[string]*CmdHandler
	*Config
}

type CmdHandler struct {
}

func (c *CLI) Parse(args []string) {
}

func Start() {
	conf, err := LoadConfig()
	if err != nil {
		fmt.Println("=== ERROR\n", err)
		os.Exit(1)
	}
	cli := CLI{make(map[string]*CmdHandler), conf}
	cli.Parse(os.Args[1:])
}
