package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
)

type Config struct {
	ConfigFile string
	ApiUser    string
	ApiKey     string
	Server     string
	Verbose    bool
}

func LoadConfig() (*Config, error) {
	conf := &Config{}

	userDir, err := getHomeDir()
	if err != nil {
		return conf, err
	}

	flag.StringVar(&conf.ConfigFile, "configFile", fmt.Sprintf("%s%c.onapp", userDir, os.PathSeparator), "Path to config file for this command")
	flag.BoolVar(&conf.Verbose, "v", false, "Verbose logging")
	flag.Parse()

	_, err = os.Stat(conf.ConfigFile)
	var merge *Config
	if err != nil {
		fmt.Printf("Config file '%s' not found, using default values\n", conf.ConfigFile)
		merge = &Config{}
	} else {
		rawConf, err := ioutil.ReadFile(conf.ConfigFile)
		if err != nil {
			return conf, errors.New(fmt.Sprintf("Error reading from %s: %s", conf.ConfigFile, err.Error()))
		}
		confFromFile := &Config{}
		err = json.Unmarshal(rawConf, confFromFile)
		if err != nil {
			return conf, errors.New(fmt.Sprintf("Error parsing %s: %s", conf.ConfigFile, err.Error()))
		}
		merge = confFromFile
	}
	merged, err := mergeConfigs(conf, merge)
	if err != nil {
		return conf, err
	}

	if merged.ApiUser == "" || merged.ApiKey == "" || merged.Server == "" {
		return merged, errors.New(fmt.Sprintf("Looks like you haven't configured a user yet, try `%s config`\n", filepath.Base(os.Args[0])))
	}

	return conf, nil
}

/* Will be fixed in Go 1.2 */
func getHomeDir() (string, error) {
	t, err := syscall.OpenCurrentProcessToken()
	defer t.Close()
	if err != nil {
		return "", err
	}
	userDir, err := t.GetUserProfileDirectory()
	if err != nil {
		return "", err
	}
	return userDir, nil
}

/* Single depth merging, prefers values in `first` over `second` */
func mergeConfigs(first *Config, second *Config) (*Config, error) {
	merged := &Config{}
	r1 := reflect.ValueOf(*first)
	r2 := reflect.ValueOf(*second)
	for i := 0; i < r1.NumField(); i++ {
		f1 := r1.Field(i)
		f2 := r2.Field(i)
		dst := reflect.ValueOf(merged).Elem().Field(i)
		switch f1.Kind() {
		case reflect.String:
			if f1.String() == "" {
				dst.SetString(f2.String())
			} else {
				dst.SetString(f1.String())
			}
		case reflect.Bool:
			if f1.Bool() || f2.Bool() {
				dst.SetBool(true)
			}
		}
	}
	return merged, nil
}
