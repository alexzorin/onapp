package cmd

import (
	"bufio"
	"code.google.com/p/go.crypto/ssh"
	"errors"
	"fmt"
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	vmCmdDescription     = "Manage virtual machines"
	vmCmdHelp            = "See subcommands for help on managing virtual machines."
	vmCmdListDescription = "List virtual machines under your account"
	vmCmdListHelp        = "\nUsage: `onapp vm list [filter]`\n" +
		"Optionally filter by field query, e.gg onapp vm list [Label=prod Hostname=.com User=1 Memory=1024]. (case sensitive)"
	vmCmdStartDescription        = "Boots a virtual machine"
	vmCmdStartHelp               = "Boots virtual machine by id: `onapp vm start <id>."
	vmCmdStopDescription         = "Stops a virtual machine"
	vmCmdStopHelp                = "Stops a virtual machine by id: `onapp vm stop <id>`."
	vmCmdRebootDescription       = "Reboots a virtual machine"
	vmCmdRebootHelp              = "Reboots a virtual machine by id: `onapp vm stop <id>`."
	vmCmdTransactionsDescription = "Lists recent transactions on a virtual machine"
	vmCmdTransactionsHelp        = "Usage: `onapp vm transactions <id> [number_to_list]`"
	vmCmdSshDescription          = "Uses SSH and the known root password to login to the machine"
	vmCmdSshHelp                 = "Usage: `onapp vm ssh <id>`, will connect on <first_ip>:22 as root with the known root password"
	vmCmdStatDescription         = "Logs into the VM via SSH and runs vmstat, printing to stdout"
	vmCmdStatHelp                = "Usage: `onapp vm stat <id>`"
)

// Base command

type vmCmd struct{}

var vmCmdHandlers = map[string]cmdHandler{
	"list":   vmCmdList{},
	"start":  vmCmdStart{},
	"stop":   vmCmdStop{},
	"reboot": vmCmdReboot{},
	"ssh":    vmCmdSsh{},
	"stat":   vmCmdStat{},
	"tx":     vmCmdTransactions{},
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
		trimmed := strings.Trim(s, " ")
		matches := pattern.FindStringSubmatch(trimmed)
		if len(matches) != 3 {
			_, err := strconv.Atoi(trimmed)
			if err == nil {
				searches = append(searches, search{"Id", trimmed})
				break
			} else {
				log.Warnf("Search query '%s' isn't valid\n", s)
			}
		} else {
			searches = append(searches, search{matches[1], matches[2]})
		}
	}
	asList := list.AsList()
	for _, s := range searches {
		asList = ctx.Search(s, asList)
	}
	log.Infof("%35.35s   #%-3s   %-5s   %-9s   %-15s   %-8s   %-11s   %-8s\n",
		"Label", "ID", "HV", "User", "First IP", "Status", "CPUs", "RAM")
	for item := asList.Front(); item != nil; item = item.Next() {
		vm := (item.Value).(onapp.VirtualMachine)
		log.Infof("%35.35s   #%-3d   HV-%-2d   User %-4d   %-15s   %-18s %5d  %10dM\n",
			vm.Label, vm.Id, vm.HV, vm.User, vm.GetIpAddress().Address, vm.BootedStringColored(), vm.Cpus, vm.Memory)
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
		vm, err := ctx.findVm(strings.Join(args, " "))
		if err != nil {
			return err
		}
		busy := ctx.checkVmBusy(vm.Id)
		if busy != nil {
			return busy
		}
		err = ctx.apiClient.VirtualMachineStartup(vm.Id)
		if err != nil {
			return err
		}
		log.Successf("Job successfully queued, waiting for boot process to start ... ")
		tx, err := ctx.awaitVmTransaction(vm.Id, "startup_virtual_machine")
		if err != nil {
			return err
		}
		log.Successf("Boot process started: #%d!\n", tx.Id)
		return nil
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
		vm, err := ctx.findVm(strings.Join(args, " "))
		if err != nil {
			return err
		}
		busy := ctx.checkVmBusy(vm.Id)
		if busy != nil {
			return busy
		}
		err = ctx.apiClient.VirtualMachineShutdown(vm.Id)
		if err != nil {
			return err
		}
		log.Successf("Job successfully queued, waiting for shutdown process to start ... ")
		tx, err := ctx.awaitVmTransaction(vm.Id, "stop_virtual_machine")
		if err != nil {
			return err
		}
		log.Successf("Shutdown process started: #%d!\n", tx.Id)
		return nil
	}
}

func (c vmCmdStop) Description() string {
	return vmCmdStopDescription
}

func (c vmCmdStop) Help(args []string) {
	log.Infoln(vmCmdStopHelp)
}

// Reboot command
type vmCmdReboot struct{}

func (c vmCmdReboot) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	} else {
		vm, err := ctx.findVm(strings.Join(args, " "))
		if err != nil {
			return err
		}
		busy := ctx.checkVmBusy(vm.Id)
		if busy != nil {
			return busy
		}
		err = ctx.apiClient.VirtualMachineReboot(vm.Id)
		if err != nil {
			return err
		}
		log.Successf("Job successfully queued, waiting for reboot process to start ... ")
		tx, err := ctx.awaitVmTransaction(vm.Id, "reboot_virtual_machine")
		if err != nil {
			return err
		}
		log.Successf("Reboot process started: #%d!\n", tx.Id)
		return nil
	}
}

func (c vmCmdReboot) Description() string {
	return vmCmdRebootDescription
}

func (c vmCmdReboot) Help(args []string) {
	log.Infoln(vmCmdRebootHelp)
}

// Transactions command
type vmCmdTransactions struct{}

func (c vmCmdTransactions) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	}
	vm, err := ctx.findVm(strings.Join(args, " "))
	if err != nil {
		return err
	}
	nList := 10
	if len(args) == 2 {
		nList, err = strconv.Atoi(strings.Trim(args[1], " "))
		if err != nil {
			return err
		}
	}
	txns, err := ctx.apiClient.VirtualMachineGetTransactions(vm.Id)
	if err != nil {
		return err
	}
	for i := 0; i <= nList && i < len(txns); i++ {
		tx := txns[i]
		t, err := tx.CreatedAtTime()
		if err != nil {
			log.Errorln(err)
			continue
		}
		log.Infof("%25.25s   #%-6d   %-25.25s   %10s\n", t, tx.Id, tx.Action, tx.StatusColored())
	}
	return nil
}

func (c vmCmdTransactions) Description() string {
	return vmCmdTransactionsDescription
}

func (c vmCmdTransactions) Help(args []string) {
	log.Infoln(vmCmdTransactionsHelp)
}

// SSH command
type vmCmdSsh struct{}

func (c vmCmdSsh) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	}
	vm, err := ctx.findVm(strings.Join(args, " "))
	if err != nil {
		return err
	}
	if !vm.Booted {
		return errors.New("Virtual machine isn't booted")
	}
	sshCmd, err := exec.LookPath("ssh")
	if err != nil {
		return err
	}
	sshArgs := []string{"ssh", fmt.Sprintf("root@%s", vm.GetIpAddress().Address)}

	log.Infof("If prompted, enter %s as the password\n", vm.RootPassword)

	err = syscall.Exec(sshCmd, sshArgs, os.Environ())
	if err != nil {
		return err
	}
	return nil
}

func (c vmCmdSsh) Description() string {
	return vmCmdSshDescription
}

func (c vmCmdSsh) Help(args []string) {
	log.Infoln(vmCmdSshHelp)
}

// Stat command
type vmCmdStat struct{}

type vmPassword string

func (s vmPassword) Password(user string) (string, error) {
	return string(s), nil
}

func (c vmCmdStat) Run(args []string, ctx *cli) error {
	if len(args) == 0 {
		c.Help(args)
		return nil
	}
	vm, err := ctx.findVm(strings.Join(args, " "))
	if err != nil {
		return err
	}
	if !vm.Booted {
		return errors.New("Virtual machine isn't booted")
	}
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.ClientAuth{
			ssh.ClientAuthPassword(vmPassword(vm.RootPassword)),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", vm.GetIpAddress().Address, 22), config)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = os.Stdout
	err = session.Run("echo && uptime && echo && free -m && echo && vmstat 1 10")
	if err != nil {
		return err
	}
	return nil
}

func (c vmCmdStat) Description() string {
	return vmCmdStatDescription
}

func (c vmCmdStat) Help(args []string) {
	log.Infoln(vmCmdStatHelp)
}

// Shared funcs

// lev distance, credit athiwatc@github, https://groups.google.com/forum/#!topic/golang-nuts/i1HMYf_DuvA
func compare(a, b string) int {
	var cost int
	d := make([][]int, len(a)+1)
	for i := 0; i < len(d); i++ {
		d[i] = make([]int, len(b)+1)
	}
	var min1, min2, min3 int
	for i := 0; i < len(d); i++ {
		d[i][0] = i
	}
	for i := 0; i < len(d[0]); i++ {
		d[0][i] = i
	}
	for i := 1; i < len(d); i++ {
		for j := 1; j < len(d[0]); j++ {
			if a[i-1] == b[j-1] {
				cost = 0
			} else {
				cost = 1
			}
			min1 = d[i-1][j] + 1
			min2 = d[i][j-1] + 1
			min3 = d[i-1][j-1] + cost
			d[i][j] = int(math.Min(math.Min(float64(min1), float64(min2)), float64(min3)))
		}
	}
	return d[len(a)][len(b)]
}

func (ctx *cli) findVm(query string) (onapp.VirtualMachine, error) {
	query = strings.Trim(strings.ToLower(query), " ")
	// Easiest option is when user passes id directly
	asInt, err := strconv.Atoi(query)
	if err == nil {
		return ctx.apiClient.GetVirtualMachine(asInt)
	}
	vms, err := ctx.apiClient.GetVirtualMachines()
	if err != nil {
		return onapp.VirtualMachine{}, err
	}
	var candidate onapp.VirtualMachine
	candidateDist := 1000
	for _, v := range vms {
		// Don't bother testing closeness if we have an exact match
		if strings.ToLower(v.Label) == query {
			return v, nil
		}
		if strings.ToLower(v.Hostname) == query {
			return v, nil
		}
		dist := compare(v.Label, query)
		if dist < candidateDist {
			candidate = v
			candidateDist = dist
		}
		dist = compare(v.Hostname, query)
		if dist < candidateDist {
			candidate = v
			candidateDist = dist
		}
	}
	if candidate.Id == 0 {
		return candidate, errors.New("Couldn't find a VM matching that")
	} else {
		log.Infof("Inexact match found for '%s': (#%d, %s) - do you want to continue? [y/n]: ", query, candidate.Id, candidate.Label)
		reader := bufio.NewReader(os.Stdin)
		resp, err := reader.ReadString('\n')
		if err != nil {
			return onapp.VirtualMachine{}, err
		}
		if strings.ToLower(resp)[0] != 'y' {
			return onapp.VirtualMachine{}, errors.New("User cancelled action")
		}
		return candidate, nil
	}
}

func (ctx *cli) awaitVmTransaction(vmId int, transType string) (onapp.Transaction, error) {
	var timeout int
	start := time.Now().Add(-5 * time.Second)
	for {
		if timeout > 30 {
			return onapp.Transaction{}, errors.New("Gave up waiting for transaction")
		}
		txns, err := ctx.apiClient.VirtualMachineGetTransactions(vmId)
		if err != nil {
			return onapp.Transaction{}, err
		}
		for i := 0; i < 5 && i < len(txns); i++ {
			tx := txns[i]
			t, err := tx.CreatedAtTime()
			// Couldn't parse time
			if err != nil {
				continue
			}
			// Job started before our job
			if t.Before(start) {
				break
			}
			// Wrong type of job
			if tx.Action != transType {
				continue
			}
			return tx, nil
		}
		<-time.After(5 * time.Second)
		timeout += 5
	}
}

func (ctx *cli) checkVmBusy(id int) error {
	busy, err := ctx.apiClient.VirtualMachineGetLatestTransaction(id, "running")
	if err != nil {
		return err
	}
	if busy.IsValid() {
		log.Warnf("This VM is currently running a transaction: %s\n", busy.Action)
		log.Warnf("Do you want to queue another action anyway? [y/n]: ")

		reader := bufio.NewReader(os.Stdin)
		resp, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.ToLower(resp)[0] == 'n' {
			return errors.New("User cancelled action")
		}
	}
	return nil
}
