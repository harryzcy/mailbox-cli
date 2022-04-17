package email

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Client struct {
	APIID       string
	Region      string
	Endpoint    string
	Credentials aws.CredentialsProvider
}

func (c Client) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/", c.APIID, c.Region)
}

var ioReadall = ioutil.ReadAll

func (c Client) request(ctx context.Context, method string, path string, query url.Values, payload []byte) (string, error) {
	body := bytes.NewReader(payload)
	req, err := http.NewRequestWithContext(ctx, method, c.getEndpoint()+path, body)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = query.Encode()

	err = SignSDKRequest(ctx, req, &SignSDKRequestOptions{
		Credentials: c.Credentials,
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

	data, err := ioReadall(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
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
	if err != nil {
		return "", err
	}
	c.Credentials = cfg.Credentials

	q := url.Values{}
	addQuery(q, "type", options.Type)
	addQuery(q, "year", options.Year)
	addQuery(q, "month", options.Month)
	addQuery(q, "order", options.Order)
	addQuery(q, "next_cursor", options.NextCursor)
	result, err := c.request(ctx, http.MethodGet, "/emails", q, nil)

	return string(result), nil
}
