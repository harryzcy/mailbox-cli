package email

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
)

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

type ListOptions struct {
	Type       string
	Year       string
	Month      string
	Order      string // asc or desc (default)
	NextCursor string
}

func (o ListOptions) check() error {
	if o.Type != EmailTypeInbox && o.Type != EmailTypeDraft && o.Type != EmailTypeSent {
		return errors.New("invalid type")
	}

	if o.Order != "" && o.Order != OrderAsc && o.Order != OrderDesc {
		return errors.New("invalid order")
	}

	return nil
}

func (c *Client) List(options ListOptions) (string, error) {
	if options.Type == "" {
		options.Type = EmailTypeInbox
	}
	if err := options.check(); err != nil {
		return "", err
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	url := c.getEndpoint() + "/emails"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	addQuery(q, "type", options.Type)
	addQuery(q, "year", options.Year)
	addQuery(q, "month", options.Month)
	addQuery(q, "order", options.Order)
	addQuery(q, "next_cursor", options.NextCursor)
	req.URL.RawQuery = q.Encode()

	err = SignSDKRequest(ctx, req, &SignSDKRequestOptions{
		Credentials: cfg.Credentials,
		Payload:     []byte(""),
		Region:      c.Region,
	})
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
