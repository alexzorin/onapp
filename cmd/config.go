package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/alexzorin/onapp"
	"github.com/alexzorin/onapp/cmd/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	configCmdDescription = "Configure this tool"
	configCmdHelp        = "This command is an interactive wizard to help you setup this tool.\n" +
		"You will need your OnApp dashboard URL, email address and API key to complete it."
)

type config struct {
	ConfigFile string
	ApiUser    string
	ApiKey     string
	Server     string
}

type configCmd struct {
}

func (c configCmd) Run(args []string, ctx *cli) error {

	log.Infoln("This is the configuration wizard. Please provide the following: \n")
	reader := bufio.NewReader(os.Stdin)

	log.Infof("Hostname of the OnApp dashboard: ")
	host, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Infof("API Username (generally the email address): ")
	user, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	log.Infof("API Key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	log.Infof("Test these details? [y/n]: ")
	doTest, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	ctx.config.Server = strings.Trim(host, "\r\n")
	ctx.config.ApiUser = strings.Trim(user, "\r\n")
	ctx.config.ApiKey = strings.Trim(apiKey, "\r\n")

	if strings.ToLower(doTest)[0] == 'y' {
		err := c.testCredentials(ctx.config.Server, ctx.config.ApiUser, ctx.config.ApiKey)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(ctx.config.ConfigFile)
	if err == nil {
		log.Infof("Config file already exists at '%s', overwrite? [y/n]: ", ctx.config.ConfigFile)
		cont, err := reader.ReadString('\n')
		if err != nil || strings.ToLower(cont)[0] != 'y' {
			return errors.New("User aborted saving configuration")
		}
	}

	err = ctx.config.save()
	if err != nil {
		log.Errorln(err)
	} else {
		log.Successln("\nSaved configuration to", ctx.config.ConfigFile)
	}

	return nil
}

func (c configCmd) testCredentials(host string, user string, pass string) error {
	client, err := onapp.NewClient(host, user, pass)
	if err != nil {
		return err
	}
	_, err = client.GetProfile()
	if err != nil {
		return err
	}
	return nil
}

func (c configCmd) Description() string {
	return configCmdDescription
}

func (c configCmd) Help(args []string) {
	log.Infoln(configCmdHelp)
}

func (c *config) save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.ConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig() (*config, error) {
	conf := &config{}

	u, err := user.Current()
	if err != nil {
		return conf, err
	}

	flag.StringVar(&conf.ConfigFile, "configFile", fmt.Sprintf("%s%c.onapp", u.HomeDir, os.PathSeparator), "Path to config file")
	flag.Parse()

	_, err = os.Stat(conf.ConfigFile)
	var merge *config
	if err != nil {
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
		log.Warnf("You haven't configured yet. Try `%s config`.\n", filepath.Base(os.Args[0]))
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
