package cmd

import (
	"os"
	"testing"
)

func TestMergeConfigs(t *testing.T) {
	c1 := &Config{ApiUser: "a", ApiKey: "abcd"}
	c2 := &Config{ApiUser: "b", Verbose: true}
	merged, err := mergeConfigs(c1, c2)
	if err != nil {
		t.Error(err)
	}
	if merged.ApiUser != "a" {
		t.Fatal("ApiUser shouldn't have been overriden")
	}
	if !merged.Verbose {
		t.Fatal("Verbose should have been true")
	}
}

func TestHomeDir(t *testing.T) {
	dir, err := getHomeDir()
	if err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Error(err)
	}
}
