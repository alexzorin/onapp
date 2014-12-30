package onapp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var cl http.Client

func (c *Client) getReq(path ...string) ([]byte, error, int) {
	url := c.makeUri(path...)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err, -1
	}
	req.SetBasicAuth(c.apiUser, c.apiPassword)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err, -1
	}
	return c.readResponse(resp)
}

func (c *Client) postReq(body string, path ...string) ([]byte, error, int) {
	url := c.makeUri(path...)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err, -1
	}
	req.SetBasicAuth(c.apiUser, c.apiPassword)
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err, -1
	}
	return c.readResponse(resp)
}

func (c *Client) readResponse(resp *http.Response) ([]byte, error, int) {
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, resp.StatusCode
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil,
			errors.New(fmt.Sprintf("Bad response on '%s' call: HTTP %d - %s\n%s", resp.Request.URL, resp.StatusCode, resp.Status, data)),
			resp.StatusCode
	}
	return data, nil, resp.StatusCode
}

// Concatenates the path elements together, and then to the dashboard URL.
func (c *Client) makeUri(toConcat ...string) string {
	buf := bytes.NewBufferString("")
	if strings.Index(c.Server, "http://") == -1 && strings.Index(c.Server, "https://") == -1 {
		buf.WriteString("https://")
	}
	buf.WriteString(c.Server)
	if c.Server[len(c.Server)-1] != '/' {
		buf.WriteByte('/')
	}
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
	User Profile
}
