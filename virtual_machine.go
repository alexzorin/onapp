package onapp

import (
	"container/list"
	"encoding/json"
	"github.com/alexzorin/onapp/cmd/log"
)

// sort.Sort'ing over this type will
// sort by UserId
type VirtualMachines []VirtualMachine

// The OnApp Virtual Machine as according to /virtual_machines.json
type VirtualMachine struct {
	Label     string `json:"label"`
	Booted    bool   `json:"booted"`
	Hostname  string `json:"hostname"`
	HV        int    `json:"hypervisor_id"`
	Cpus      int    `json:"cpus"`
	CpuShares int    `json:"cpu_shares"`
	Memory    int    `json:"memory"`
	Template  string `json:"template_label"`
	User      int    `json:"user_id"`
}

func (vm *VirtualMachine) BootedString() string {
	if vm.Booted {
		return "Booted"
	} else {
		return "Offline"
	}
}

func (vm *VirtualMachine) BootedStringColored() string {
	if vm.Booted {
		return log.ColorString("Booted", log.GREEN)
	} else {
		return log.ColorString("Offline", log.RED)
	}
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

// Fetches a list of Virtual Machines from the dashboard server
func (c *Client) GetVirtualMachines() (VirtualMachines, error) {
	data, err := c.getReq("virtual_machines.json")
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
	}
	return vms, nil
}
