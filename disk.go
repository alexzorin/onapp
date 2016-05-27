package onapp

import (
	"encoding/json"
	"strconv"
)

type Disks []Disk

type Disk struct {
	AddToFreebsdFstab              interface{} `json:"add_to_freebsd_fstab"`
	AddToLinuxFstab                interface{} `json:"add_to_linux_fstab"`
	Built                          bool        `json:"built"`
	BurstBw                        int         `json:"burst_bw"`
	BurstIops                      int         `json:"burst_iops"`
	CreatedAt                      string      `json:"created_at"`
	DataStoreID                    int         `json:"data_store_id"`
	DiskSize                       int         `json:"disk_size"`
	DiskVMNumber                   int         `json:"disk_vm_number"`
	FileSystem                     string      `json:"file_system"`
	ID                             int         `json:"id"`
	Identifier                     string      `json:"identifier"`
	IntegratedStorageCacheEnabled  bool        `json:"integrated_storage_cache_enabled"`
	IntegratedStorageCacheOverride bool        `json:"integrated_storage_cache_override"`
	IntegratedStorageCacheSettings struct {
	} `json:"integrated_storage_cache_settings"`
	Iqn              interface{} `json:"iqn"`
	IsSwap           bool        `json:"is_swap"`
	Label            string      `json:"label"`
	Locked           bool        `json:"locked"`
	MaxBw            int         `json:"max_bw"`
	MaxIops          int         `json:"max_iops"`
	MinIops          int         `json:"min_iops"`
	MountPoint       interface{} `json:"mount_point"`
	Primary          bool        `json:"primary"`
	UpdatedAt        string      `json:"updated_at"`
	VirtualMachineID int         `json:"virtual_machine_id"`
	VolumeID         interface{} `json:"volume_id"`
	HasAutobackups   bool        `json:"has_autobackups"`
}

type DiskSchedules []DiskSchedule

type DiskSchedule struct {
	Action         string            `json:"action"`
	CreatedAt      string            `json:"created_at"`
	Duration       int               `json:"duration"`
	FailureCount   int               `json:"failure_count"`
	ID             int               `json:"id"`
	Params         interface{}       `json:"params"`
	Period         string            `json:"period"`
	RotationPeriod int               `json:"rotation_period"`
	StartAt        string            `json:"start_at"`
	Status         string            `json:"status"`
	TargetID       int               `json:"target_id"`
	TargetType     string            `json:"target_type"`
	UpdatedAt      string            `json:"updated_at"`
	UserID         int               `json:"user_id"`
	ScheduleLogs   []DiskScheduleLog `json:"schedule_logs"`
}

type DiskScheduleLog struct {
	Log struct {
		CreatedAt  string `json:"created_at"`
		ID         int    `json:"id"`
		LogOutput  string `json:"log_output"`
		ScheduleID int    `json:"schedule_id"`
		Status     string `json:"status"`
		UpdatedAt  string `json:"updated_at"`
	} `json:"schedule_log"`
}

func (c *Client) GetVirtualMachineDisks(vmId int) (Disks, error) {
	data, err, _ := c.getReq("/virtual_machines/", strconv.Itoa(vmId), "/disks.json")
	if err != nil {
		return Disks{}, err
	}
	var out []map[string]Disk
	err = json.Unmarshal(data, &out)
	if err != nil {
		return Disks{}, err
	}

	backups := make([]Disk, len(out))
	for i := range backups {
		backups[i] = out[i]["disk"]
	}
	return backups, nil
}

func (c *Client) GetVirtualMachineDiskSchedules(vmId, diskId int) (DiskSchedules, error) {
	data, err, _ := c.getReq(
		"/virtual_machines/", strconv.Itoa(vmId), "/disks/",
		strconv.Itoa(diskId), "/schedules.json")
	if err != nil {
		return DiskSchedules{}, err
	}
	var out []map[string]DiskSchedule
	err = json.Unmarshal(data, &out)
	if err != nil {
		return DiskSchedules{}, err
	}

	ds := make([]DiskSchedule, len(out))
	for i := range ds {
		ds[i] = out[i]["schedule"]
	}
	return ds, nil
}
