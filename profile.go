package onapp

import (
	"encoding/json"
)

// The OnApp user profile as according to /profile.json
type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Login     string `json:"login"`
	Id        int    `json:"id"`
	Email     string `json:"email"`
}

// Fetches the user profile from the dashboard server
func (c *Client) GetProfile() (*Profile, error) {
	data, err := c.getReq("profile.json")
	if err != nil {
		return nil, err
	}
	var j jsonOuterFields
	err = json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	return &j.User, nil
}
