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
	if s != "http://dashboard.example.org/this/is/atest" {
		t.Fail()
	}
}

func TestUnmarshalInner(t *testing.T) {
	testJson := `{"outer": { "inner": 123 }}`
	c := newTestClient()
	inner, err := c.unmarshalInner([]byte(testJson), "Outer")
	if err != nil {
		t.Error(err)
		return
	}
	if inner == nil {
		t.Error("Inner element was nil")
		return
	}
	asserted, ok := inner.(testStruct)
	if !ok {
		t.Error("Couldn't cast to our test struct")
	}
	if asserted.Inner != 123 {
		t.Error("Asserted struct's value was wrong")
	}
}
