package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
)

const (
	configCmdDescription = "Configure this tool"
)

type config struct {
	ConfigFile string
	ApiUser    string
	ApiKey     string
	Server     string
	Verbose    bool
}

type configCmd struct {
}

func (c configCmd) Run(args []string, ctx *cli) error {
	return nil
}

func (c configCmd) Description() string {
	return configCmdDescription
}

func loadConfig() (*config, error) {
	conf := &config{}

	u, err := user.Current()
	if err != nil {
		return conf, err
	}

	flag.StringVar(&conf.ConfigFile, "configFile", fmt.Sprintf("%s%c.onapp", u.HomeDir, os.PathSeparator), "Path to config file")
	flag.BoolVar(&conf.Verbose, "v", false, "Verbose logging")
	flag.Parse()

	_, err = os.Stat(conf.ConfigFile)
	var merge *config
	if err != nil {
		warnf("Config file '%s' not found\n", conf.ConfigFile)
		merge = &config{}
	} else {
		rawConf, err := ioutil.ReadFile(conf.ConfigFile)
		if err != nil {
			return conf, errors.New(fmt.Sprintf("Error reading from %s: %s", conf.ConfigFile, err.Error()))
		}
		confFromFile := &config{}
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
		warnf("You haven't configured yet: `%s config`\n", filepath.Base(os.Args[0]))
	}

	return merged, nil
}

/* Single depth merging, prefers values in `first` over `second` */
func mergeConfigs(first *config, second *config) (*config, error) {
	merged := &config{}
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
