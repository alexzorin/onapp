package cmd

import (
	"encoding/json"
	"errors"
	"github.com/alexzorin/onapp"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const (
	cacheFileName = ".onapp_cache"
)

var (
	ErrCacheDoesntExist = errors.New("Cache doesn't exist")
)

type Cache interface {
	GetVirtualMachines() (onapp.VirtualMachines, error)
	Clear()
	Store(onapp.VirtualMachines) error
}

type fileBackedCache struct {
}

func (c *fileBackedCache) GetVirtualMachines() (onapp.VirtualMachines, error) {
	p, err := c.getCacheFile(false)
	if err != nil {
		return nil, ErrCacheDoesntExist
	}
	cf, err := os.Open(p)
	if err != nil {
		return nil, errors.New("Cache unopenable - " + err.Error())
	}
	defer cf.Close()
	data, err := ioutil.ReadAll(cf)
	if err != nil {
		return nil, errors.New("Couldn't read the cache file - " + err.Error())
	}
	var out onapp.VirtualMachines
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, errors.New("Couldn't parse the cache file - " + err.Error())
	}
	return out, nil
}

func (c *fileBackedCache) Store(vms onapp.VirtualMachines) error {
	// Clean the vm info
	for k, _ := range vms {
		vm := &vms[k]
		vm.VncPassword = ""
		vm.RootPassword = ""
	}
	p, err := c.getCacheFile(true)
	if err != nil {
		return err
	}
	data, err := json.Marshal(vms)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(p, data, 0644); err != nil {
		return err
	}
	return nil
}

func (c *fileBackedCache) Clear() {
	c.getCacheFile(true)

}

func (c *fileBackedCache) getCacheFile(unlink bool) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	path := u.HomeDir + string(filepath.Separator) + cacheFileName
	_, err = os.Stat(path)
	if err != nil && !unlink {
		return "", err
	} else if err == nil && unlink {
		os.Remove(path)
	}
	return path, nil
}
