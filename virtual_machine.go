package onapp

import (
	"container/list"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/alexzorin/onapp/log"
	"github.com/alexzorin/rfc5735"
)

// sort.Sort'ing over this type will
// sort by UserId
type VirtualMachines []VirtualMachine

// The OnApp Virtual Machine as according to /virtual_machines.json
type VirtualMachine struct {
	client         *Client
	Id             int                    `json:"id"`
	Label          string                 `json:"label"`
	Booted         bool                   `json:"booted"`
	Hostname       string                 `json:"hostname"`
	HV             int                    `json:"hypervisor_id"`
	Cpus           int                    `json:"cpus"`
	CpuShares      int                    `json:"cpu_shares"`
	Memory         int                    `json:"memory"`
	Template       string                 `json:"template_label"`
	User           int                    `json:"user_id"`
	Locked         bool                   `json:"locked"`
	RootPassword   string                 `json:"initial_root_password"`
	IpAddressesRaw []map[string]IpAddress `json:"ip_addresses"`
	VncPassword    string                 `json:"remote_access_password"`
	AdminNote      string                 `json:"admin_note"`
}

// IP address of a virtual machine as represented by /virtual_machines/:id.json
type IpAddress struct {
	Address        string `json:"address"`
	Gateway        string `json:"gateway"`
	Broadcast      string `json:"broadcast"`
	NetworkAddress string `json:"network_address"`
	Netmask        string `json:"netmask"`
}

// Remote Access Session
type RemoteAccessSession struct {
	Port int `json:"port"`
}

func (c *Client) GetRemoteAccessSession(id int) (RemoteAccessSession, error) {
	data, err, _ := c.getReq("/virtual_machines/", strconv.Itoa(id), "/console.json")
	if err != nil {
		return RemoteAccessSession{}, err
	}
	var out map[string]RemoteAccessSession
	err = json.Unmarshal(data, &out)
	if err != nil {
		return RemoteAccessSession{}, err
	} else {
		return out["remote_access_session"], err
	}
}

// Fetches a list of Virtual Machines from the dashboard server
func (c *Client) GetVirtualMachines() (VirtualMachines, error) {
	data, err, _ := c.getReq("virtual_machines.json")
	if err != nil {
		return nil, err
	}
	var out []map[string]VirtualMachine
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	vms := make([]VirtualMachine, len(out))
	for i := range vms {
		vms[i] = out[i]["virtual_machine"]
		vms[i].client = c
	}
	return vms, nil
}

// Fetches an individual Virtual Machine from the dashboard server
func (c *Client) GetVirtualMachine(id int) (VirtualMachine, error) {
	data, err, _ := c.getReq("virtual_machines/", strconv.Itoa(id), ".json")
	if err != nil {
		return VirtualMachine{}, err
	}
	var out map[string]VirtualMachine
	err = json.Unmarshal(data, &out)
	if err != nil {
		return VirtualMachine{}, err
	}
	vm := out["virtual_machine"]
	vm.client = c
	return vm, nil
}

func (c *Client) VirtualMachineStartup(id int) error {
	_, err, rc := c.postReq("", "virtual_machines/", strconv.Itoa(id), "/startup.json")
	if rc == 422 {
		return errors.New("HTTP 422 - VM can't currently be booted")
	}
	return err
}

func (c *Client) VirtualMachineShutdown(id int) error {
	_, err, rc := c.postReq("", "virtual_machines/", strconv.Itoa(id), "/shutdown.json")
	if rc == 422 {
		return errors.New("HTTP 422 - VM can't currently be shut down")
	}
	return err
}

func (c *Client) VirtualMachineReboot(id int) error {
	_, err, rc := c.postReq("", "virtual_machines/", strconv.Itoa(id), "/reboot.json")
	if rc == 422 {
		return errors.New("HTTP 422 - VM can't currently be rebooted")
	}
	return err
}

func (c *Client) VirtualMachineGetTransactions(vmId int) (Transactions, error) {
	return c.getTransactions(vmId)
}

func (c *Client) VirtualMachineGetLatestTransaction(vmId int, statuses ...string) (Transaction, error) {
	txns, err := c.VirtualMachineGetTransactions(vmId)
	if err != nil {
		return Transaction{}, err
	}
	if len(statuses) == 0 && len(txns) > 0 {
		return txns[0], nil
	}
	for _, t := range txns {
		if len(statuses) > 0 {
			for _, s := range statuses {
				if t.Status == s {
					return t, nil
				}
			}
		}
	}
	return Transaction{}, nil
}

func (vm *VirtualMachine) Startup() error {
	return vm.client.VirtualMachineStartup(vm.Id)
}

func (vm *VirtualMachine) Shutdown() error {
	return vm.client.VirtualMachineShutdown(vm.Id)
}

func (vm *VirtualMachine) Reboot() error {
	return vm.client.VirtualMachineReboot(vm.Id)
}

func (vm *VirtualMachine) GetTransactions() (Transactions, error) {
	return vm.client.VirtualMachineGetTransactions(vm.Id)
}

func (vm *VirtualMachine) GetRunningTransaction() (Transaction, error) {
	return vm.client.VirtualMachineGetLatestTransaction(vm.Id, "running")
}

func (vm *VirtualMachine) GetIpAddresses() ([]IpAddress, error) {
	var addrs []IpAddress
	for _, v := range vm.IpAddressesRaw {
		addrs = append(addrs, v["ip_address"])
	}
	return addrs, nil
}

// This prefers non-rfc5735 addresses
// if possible
func (vm *VirtualMachine) GetIpAddress() IpAddress {
	var out IpAddress
	ips, err := vm.GetIpAddresses()
	if err != nil {
		return IpAddress{}
	} else {
		if len(ips) > 0 {
			out = ips[0]
		}
		if len(ips) > 1 {
			for _, i := range ips {
				if !rfc5735.IsReservedString(i.Address) {
					out = i
					break
				}
			}
		}
	}
	return out
}

func (vm *VirtualMachine) BootedString() string {
	if vm.Locked {
		return "Locked"
	} else if vm.Booted {
		return "Booted"
	} else {
		return "Offline"
	}
}

func (vm *VirtualMachine) BootedStringColored() string {
	if vm.Locked {
		return log.ColorString("LOCKED", log.YELLOW)
	} else if vm.Booted {
		return log.ColorString("Booted", log.GREEN)
	} else {
		return log.ColorString("Offline", log.RED)
	}
}

func (vm *VirtualMachine) GetRemoteAccessSession() (*RemoteAccessSession, error) {
	ras, err := vm.client.GetRemoteAccessSession(vm.Id)
	if err != nil {
		return &RemoteAccessSession{}, err
	}
	newVm, err := vm.client.GetVirtualMachine(vm.Id)
	if err != nil {
		return &RemoteAccessSession{}, err
	}
	*vm = newVm
	return &ras, nil
}

func (vms VirtualMachines) AsList() list.List {
	var l list.List
	for _, v := range vms {
		l.PushBack(v)
	}
	return l
}

func (vms VirtualMachines) Swap(i, j int) {
	vms[i], vms[j] = vms[j], vms[i]
}

func (vms VirtualMachines) Len() int {
	return len(vms)
}

func (vms VirtualMachines) Less(i, j int) bool {
	return vms[i].User < vms[j].User
}
