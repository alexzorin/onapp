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
		fmt.Println("No command passed, usage is as follows:")
		printUsage()
		return
	}
	if handler, ok := cmdHandlers[args[0]]; ok {
		err := handler.Run(args[1:], c)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("%s is an unknown command\n", args[0])
		printUsage()
	}
}

func Start() {
	conf, err := LoadConfig()
	if err != nil {
		fmt.Println("=== ERROR\n", err)
		os.Exit(1)
	}
	cli := CLI{conf}

	cli.parse(cleanArgs(os.Args[1:]))
}

func printUsage() {
	fmt.Println("============ USAGE ============")
	fmt.Println("Available commands:")
	for k, v := range cmdHandlers {
		fmt.Printf("  %s - %s\n", k, v.Description())
	}
	fmt.Println("General options:")
	flag.PrintDefaults()
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
