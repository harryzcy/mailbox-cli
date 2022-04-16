package email

import "fmt"

type Client struct {
	APIID    string
	Region   string
	Endpoint string
}

func (c Client) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/", c.APIID, c.Region)
}
