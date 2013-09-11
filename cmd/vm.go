package cmd

import (
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/cmd/log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	vmCmdDescription     = "Manage virtual machines"
	vmCmdHelp            = "See subcommands for help on managing virtual machines."
	vmCmdListDescription = "List virtual machines under your account"
	vmCmdListHelp        = "\nUsage: `onapp vm list [filter]`\n" +
		"Optionally filter by field query, e.gg onapp vm list Label=prod [Hostname=.com User=1 Memory=1024]. (case sensitive)"
	vmCmdStartDescription = "Boots a virtual machine"
	vmCmdStartHelp        = "Boots virtual machine by id: `onapp vm start 1`."
	vmCmdStopDescription  = "Stops a virtual machine"
	vmCmdStopHelp         = "Stops a virtual machine by id: `onapp vm stop 1`."
)

// Base command

type vmCmd struct{}

var vmCmdHandlers = map[string]cmdHandler{
	"list":  vmCmdList{},
	"start": vmCmdStart{},
	"stop":  vmCmdStop{},
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
	list, err := ctx.apiClient.GetVirtualMachines()
	if err != nil {
		return err
	}
	sort.Sort(list)
	var searches []search
	pattern := regexp.MustCompile("^(\\w+)=(\\w+)$")
	for _, s := range args {
		matches := pattern.FindStringSubmatch(strings.Trim(s, " "))
		if len(matches) != 3 {
			log.Warnf("Search query '%s' isn't valid\n", s)
		} else {
			searches = append(searches, search{matches[1], matches[2]})
		}
	}
	asList := list.AsList()
	for _, s := range searches {
		asList = ctx.Search(s, asList)
	}
	for item := asList.Front(); item != nil; item = item.Next() {
		vm := (item.Value).(onapp.VirtualMachine)
		log.Infof("%25.25s   #%-3d   HV-%-2d   User %-4d   %-18s %2d CPUs  %6dM RAM   %-30.25s\n",
			vm.Label, vm.Id, vm.HV, vm.User, vm.BootedStringColored(), vm.Cpus, vm.Memory, vm.Template)
	}
	return nil
}

func (c vmCmdList) Description() string {
	return vmCmdListDescription
}

func (c vmCmdList) Help(args []string) {
	log.Infoln(vmCmdListHelp)
	log.Infoln("\nField names are as follows: ")
	log.Infof("%+v\n\n", &onapp.VirtualMachine{})
}

// Start command
type vmCmdStart struct{}

func (c vmCmdStart) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	} else {
		id, err := strconv.Atoi(strings.Trim(args[0], " "))
		if err != nil {
			return err
		}
		return ctx.apiClient.VirtualMachineStartup(id)
	}
}

func (c vmCmdStart) Description() string {
	return vmCmdStartDescription
}

func (c vmCmdStart) Help(args []string) {
	log.Infoln(vmCmdStartHelp)
}

// Stop command
type vmCmdStop struct{}

func (c vmCmdStop) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	} else {
		id, err := strconv.Atoi(strings.Trim(args[0], " "))
		if err != nil {
			return err
		}
		return ctx.apiClient.VirtualMachineShutdown(id)
	}
}

func (c vmCmdStop) Description() string {
	return vmCmdStopDescription
}

func (c vmCmdStop) Help(args []string) {
	log.Infoln(vmCmdStopHelp)
}
