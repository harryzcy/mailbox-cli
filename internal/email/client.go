package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

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
	c.Endpoint = fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com", c.APIID, c.Region)

	if c.Verbose {
		fmt.Printf("[DEBUG] Generated endpoint: %s\n", c.Endpoint)
	}
	return c.Endpoint
}

func (c *Client) loadCredentials(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	c.Credentials = cfg.Credentials
	return nil
}

var ioReadall = io.ReadAll

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
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Add("Content-Type", "application/json")
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Request URL: %s\n", req.URL.String())
	}

	req.Header.Set("Accept", "application/json")

	if c.Verbose {
		fmt.Printf("[DEBUG] Signing request\n")
	}

	err = SignSDKRequest(ctx, req, &SignSDKRequestOptions{
		Credentials: c.Credentials,
		Payload:     payload,
		Region:      c.Region,
		Verbose:     c.Verbose,
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
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	addQuery(q, "type", options.Type)
	addQuery(q, "year", options.Year)
	addQuery(q, "month", options.Month)
	addQuery(q, "order", options.Order)
	addQuery(q, "next_cursor", options.NextCursor)
	result, err := c.request(ctx, http.MethodGet, "/emails", q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
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
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodGet, "/emails/"+options.MessageID, q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
}

type TrashOptions struct {
	MessageID string
}

func (o TrashOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	return nil
}

func (c *Client) Trash(options TrashOptions) (string, error) {
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Trashing email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodPost, "/emails/"+options.MessageID+"/trash", q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
}

type UntrashOptions struct {
	MessageID string
}

func (o UntrashOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	return nil
}

func (c *Client) Untrash(options UntrashOptions) (string, error) {
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Untrashing email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodPost, "/emails/"+options.MessageID+"/untrash", q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
}

type DeleteOptions struct {
	MessageID string
}

func (o DeleteOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	return nil
}

func (c *Client) Delete(options DeleteOptions) (string, error) {
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Deleteing email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodDelete, "/emails/"+options.MessageID, q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
}

type CreateOptions struct {
	Subject      string   `json:"subject"`
	From         []string `json:"from"`
	To           []string `json:"to"`
	Cc           []string `json:"cc"`
	Bcc          []string `json:"bcc"`
	ReplyTo      []string `json:"replyTo"`
	Text         string   `json:"text"`
	HTML         string   `json:"html"`
	GenerateText string   `json:"generateText"`
	Send         bool     `json:"send"`

	File string `json:"-"`
}

func (o CreateOptions) check() error {
	if o.GenerateText != "on" && o.GenerateText != "off" && o.GenerateText != "auto" {
		return errors.New("invalid generate-text")
	}

	return nil
}

func (o *CreateOptions) loadFile() error {
	if o.File == "" {
		return nil
	}

	data, err := os.ReadFile(o.File)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, o)

	return err
}

func (c *Client) Create(options CreateOptions) (string, error) {
	if err := options.loadFile(); err != nil {
		return "", err
	}

	if options.GenerateText == "" {
		options.GenerateText = "auto"
	}
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Creating email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodPost, "/emails", q, body)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, err
}

type SaveOptions struct {
	MessageID    string   `json:"-"`
	Subject      string   `json:"subject"`
	From         []string `json:"from"`
	To           []string `json:"to"`
	Cc           []string `json:"cc"`
	Bcc          []string `json:"bcc"`
	ReplyTo      []string `json:"replyTo"`
	Body         string   `json:"body"`
	Text         string   `json:"text"`
	HTML         string   `json:"html"`
	GenerateText string   `json:"generateText"`
	Send         bool     `json:"send"`

	File string `json:"-"`
}

func (o SaveOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	if o.Subject == "" {
		return errors.New("invalid subject")
	}

	if len(o.From) == 0 {
		return errors.New("invalid from")
	}

	if len(o.To) == 0 {
		return errors.New("invalid to")
	}

	if o.GenerateText != "on" && o.GenerateText != "off" && o.GenerateText != "auto" {
		return errors.New("invalid generate-text")
	}

	return nil
}

func (o *SaveOptions) loadFile() error {
	if o.File == "" {
		return nil
	}

	data, err := os.ReadFile(o.File)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, o)

	return err
}

func (c *Client) Save(options SaveOptions) (string, error) {
	if err := options.loadFile(); err != nil {
		return "", err
	}

	if options.GenerateText == "" {
		options.GenerateText = "auto"
	}
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Saving email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodPut, "/emails/"+options.MessageID, q, body)

	return result, err
}

type SendOptions struct {
	MessageID string
}

func (o SendOptions) check() error {
	if o.MessageID == "" {
		return errors.New("invalid message id")
	}

	return nil
}

func (c *Client) Send(options SendOptions) (string, error) {
	if err := options.check(); err != nil {
		return "", err
	}

	if c.Verbose {
		fmt.Printf("[DEBUG] Sending email\n")
	}

	ctx := context.Background()
	err := c.loadCredentials(ctx)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	result, err := c.request(ctx, http.MethodPost, "/emails/"+options.MessageID+"/send", q, nil)
	if err != nil {
		if c.Verbose {
			fmt.Printf("[DEBUG] Error: %s\n", err)
		}
		return "", err
	}

	return result, nil
}
