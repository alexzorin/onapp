package cmd

import (
	"testing"
)

func TestMergeConfigs(t *testing.T) {
	c1 := &config{ApiUser: "a", ApiKey: "abcd"}
	c2 := &config{ApiUser: "b"}
	merged, err := mergeConfigs(c1, c2)
	if err != nil {
		t.Error(err)
	}
	if merged.ApiUser != "a" {
		t.Fatal("ApiUser shouldn't have been overriden")
	}
}
