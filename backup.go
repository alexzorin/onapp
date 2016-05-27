package onapp

import (
	"encoding/json"
	"strconv"
)

type Backups []Backup

type Backup struct {
	AllowResizeWithoutReboot bool        `json:"allow_resize_without_reboot"`
	AllowedHotMigrate        bool        `json:"allowed_hot_migrate"`
	AllowedSwap              bool        `json:"allowed_swap"`
	BackupServerID           interface{} `json:"backup_server_id"`
	BackupSize               int         `json:"backup_size"`
	Built                    bool        `json:"built"`
	BuiltAt                  string      `json:"built_at"`
	CreatedAt                string      `json:"created_at"`
	DataStoreType            string      `json:"data_store_type"`
	ID                       int         `json:"id"`
	Identifier               string      `json:"identifier"`
	Initiated                string      `json:"initiated"`
	Iqn                      interface{} `json:"iqn"`
	Locked                   bool        `json:"locked"`
	MarkedForDelete          bool        `json:"marked_for_delete"`
	MinDiskSize              int         `json:"min_disk_size"`
	MinMemorySize            int         `json:"min_memory_size"`
	Note                     string      `json:"note"`
	OperatingSystem          string      `json:"operating_system"`
	OperatingSystemDistro    string      `json:"operating_system_distro"`
	TargetID                 int         `json:"target_id"`
	TargetType               string      `json:"target_type"`
	TemplateID               int         `json:"template_id"`
	UpdatedAt                string      `json:"updated_at"`
	UserID                   int         `json:"user_id"`
	VolumeID                 interface{} `json:"volume_id"`
	BackupType               string      `json:"backup_type"`
	DiskID                   int         `json:"disk_id"`
}

func (c *Client) GetVirtualMachineBackups(vmId int) (Backups, error) {
	data, err, _ := c.getReq("/virtual_machines/", strconv.Itoa(vmId), "/backups.json")
	if err != nil {
		return Backups{}, err
	}
	var out []map[string]Backup
	err = json.Unmarshal(data, &out)
	if err != nil {
		return Backups{}, err
	}

	backups := make([]Backup, len(out))
	for i := range backups {
		backups[i] = out[i]["backup"]
	}
	return backups, nil
}
