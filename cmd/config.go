package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

type Config struct {
	ConfigFile string
	ApiUser    string
	ApiKey     string
	Server     string
}

func LoadConfig() (*Config, error) {
	conf := &Config{}

	userDir, err := getHomeDir()
	if err != nil {
		return conf, err
	}

	flag.StringVar(&conf.ConfigFile, "configFile", fmt.Sprintf("%s%c.onapp", userDir, os.PathSeparator), "Path to config file for this command")
	flag.Parse()

	_, err = os.Stat(conf.ConfigFile)
	if err != nil {
		fmt.Printf("Config file '%s' doesn't exist ...\n", conf.ConfigFile)
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
