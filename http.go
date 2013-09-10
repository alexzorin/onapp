package onapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

var cl http.Client

func (c *Client) getReq(path ...string) ([]byte, error) {
	url := c.makeUri(path...)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiUser, c.apiPassword)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Bad response on '%s' call: HTTP %d - %s", url, resp.StatusCode, resp.Status))
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Concatenates the path elements together, and then to the dashboard URL.
func (c *Client) makeUri(toConcat ...string) string {
	buf := bytes.NewBufferString("http://")
	buf.WriteString(c.Server)
	buf.WriteByte('/')
	for _, v := range toConcat {
		buf.WriteString(v)
	}
	return buf.String()
}

// This is a union of all of the possible outer fields
// returned by the OnApp API, so that we don't need to
// have pointless nesting.
// * TestUnmarshalInner expects wrap by "outer":{}
// * /profile.json wrapped by "user":{}
type jsonOuterFields struct {
	Outer interface{}
	User  Profile
}

// Since OnApp pointlessly wraps many of their API responses, we sometimes only want to
// unmarshal from an inner element. This prevents our structs from needless depth.
// This will unmarshal from a named top-level element.
func (c *Client) unmarshalInner(raw []byte, name string) (interface{}, error) {
	var iface jsonOuterFields
	err := json.Unmarshal(raw, &iface)
	if err != nil {
		return nil, err
	}
	rv := reflect.ValueOf(iface)
	outerv := rv.FieldByName(name)
	if !outerv.IsValid() {
		return nil, errors.New(fmt.Sprintf("Outer field %s was not found", name))
	}
	return outerv.Interface(), nil
}
