package onapp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type Client struct {
	Server      string
	apiUser     string
	apiPassword string
}

// Creates a new API client with the specified hostname, email address and API key.
// The hostname needs to be the DNS-resolvable hostname of the dashboard server (such as dashboard.example.org).
func NewClient(hostname string, email string, apiKey string) (*Client, error) {
	if hostname == "" || email == "" || apiKey == "" {
		return nil, errors.New("Invalid parameters to NewClient")
	}

	cl := &Client{hostname, email, apiKey}
	return cl, nil
}

// Creates a new API client using the file (defaults to ~/.onapp) or the OS environment variables
// Environment variables take precedence.
func NewClientFromSystem(file string) (*Client, error) {
	var out struct {
		Server  string `json:"Server"`
		ApiUser string `json:"ApiUser"`
		ApiKey  string `json:"ApiKey"`
	}

	us, err := user.Current()
	if err != nil {
		return nil, err
	}

	if file == "" {
		file = filepath.Join(us.HomeDir, ".onapp")
	}

	_, err = os.Stat(file)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err == nil {
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(buf, &out); err != nil {
			return nil, err
		}
	}

	h := os.Getenv("ONAPP_HOST")
	u := os.Getenv("ONAPP_USER")
	p := os.Getenv("ONAPP_PASSWORD")

	if h != "" {
		out.Server = h
	}
	if u != "" {
		out.ApiUser = u
	}
	if p != "" {
		out.ApiKey = p
	}

	return NewClient(out.Server, out.ApiUser, out.ApiKey)
}
