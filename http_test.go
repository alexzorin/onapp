package onapp

import (
	"testing"
)

func newTestClient() *Client {
	c, _ := NewClient("dashboard.example.org", "user@example.org", "1234")
	return c
}

func TestNewClient(t *testing.T) {
	_, err := NewClient("dashboard.example.org", "user@example.org", "1234")
	if err != nil {
		t.Error(err)
	}
}

func TestMakeUri(t *testing.T) {
	c := newTestClient()
	s := c.makeUri("this/", "is/", "a", "test")
	if s != "https://dashboard.example.org/this/is/atest" {
		t.Fail()
	}

	c.Server = "http://dashboard.example.org/"
	s = c.makeUri("this/", "is/", "a", "test")
	if s != "http://dashboard.example.org/this/is/atest" {
		t.Fail()
	}
}
