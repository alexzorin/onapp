package onapp

type Client struct {
	Server      string
	apiUser     string
	apiPassword string
}

// Creates a new API client with the specified hostname, email address and API key.
// The hostname needs to be the DNS-resolvable hostname of the dashboard server (such as dashboard.example.org).
func NewClient(hostname string, email string, apiKey string) (*Client, error) {
	cl := &Client{hostname, email, apiKey}
	return cl, nil
}
