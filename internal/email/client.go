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
	Verbose     bool
}

func (c *Client) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	c.Endpoint = fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/", c.APIID, c.Region)

	if c.Verbose {
		fmt.Printf("[DEBUG] Generated endpoint: %s\n", c.Endpoint)
	}
	return c.Endpoint
}

var ioReadall = ioutil.ReadAll

func (c Client) request(ctx context.Context, method string, path string, query url.Values, payload []byte) (string, error) {
	body := bytes.NewReader(payload)

	if c.Verbose {
		fmt.Printf("[DEBUG] New request: %s %s\n", method, c.getEndpoint()+path)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.getEndpoint()+path, body)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = query.Encode()

	if c.Verbose {
		fmt.Printf("[DEBUG] Request URL: %s\n", req.URL.String())
	}

	req.Header.Set("Accept", "application/json")

	if c.Verbose {
		fmt.Printf("[DEBUG] Signing request\n")
	}

	err = SignSDKRequest(ctx, req, &SignSDKRequestOptions{
		Credentials: c.Credentials,
		Payload:     []byte(""),
		Region:      c.Region,
	})
	if err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Sending request\n")
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if c.Verbose {
		fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
	}

	data, err := ioReadall(resp.Body)
	if err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Content-Type: %s\n", resp.Header["Content-Type"][0])
		fmt.Printf("[DEBUG] Response: %s\n", string(data))
	}

	if resp.Header["Content-Type"][0] == "application/json" {
		if c.Verbose {
			fmt.Printf("[DEBUG] Prettifying output\n")
		}
		return prettyResult(data)
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
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Listing emails\n")
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

type GetOptions struct {
	MessageID string
}

func (o GetOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	return nil
}

func (c *Client) Get(options GetOptions) (string, error) {
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Getting email\n")
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	c.Credentials = cfg.Credentials

	q := url.Values{}
	result, err := c.request(ctx, http.MethodGet, "/emails/"+options.MessageID, q, nil)

	return string(result), nil
}
