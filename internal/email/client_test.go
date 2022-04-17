package email

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestGetEndpoint(t *testing.T) {
	tests := []struct {
		client   Client
		endpoint string
	}{
		{
			client: Client{
				Endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com/",
			},
			endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com/",
		},
		{
			client: Client{
				APIID:  "api_id",
				Region: "us-west-2",
			},
			endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com/",
		},
	}

	for _, test := range tests {
		endpoint := test.client.getEndpoint()
		assert.Equal(t, test.endpoint, endpoint)
	}
}

func TestClient_Request(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Nanosecond) // used by one test case
	defer cancel()
	defer func() {
		ioReadall = ioutil.ReadAll
	}()

	tests := []struct {
		ctx       context.Context
		client    Client
		method    string
		path      string
		query     url.Values
		payload   []byte
		ioReadall func(io.Reader) ([]byte, error)
		err       error
	}{
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: "https://httpbin.org",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			method:    http.MethodGet,
			path:      "/get",
			query:     url.Values{},
			payload:   []byte(""),
			ioReadall: ioutil.ReadAll,
			err:       nil,
		},
		{
			client: Client{
				Endpoint: "https://httpbin.org",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			err: errors.New("net/http: nil Context"),
		},
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: "https://httpbin.org",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, errors.New("error")
				}),
			},
			err: errors.New("error"),
		},
		{
			ctx: timeoutCtx,
			client: Client{
				Endpoint: "https://httpbin.org",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			path:    "/get",
			query:   url.Values{},
			payload: []byte(""),
			err:     &url.Error{Op: "Get", URL: "https://httpbin.org/get", Err: context.DeadlineExceeded},
		},
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: "https://httpbin.org",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			path: "/get",
			ioReadall: func(r io.Reader) ([]byte, error) {
				return nil, errors.New("error")
			},
			err: errors.New("error"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.ioReadall != nil {
				ioReadall = test.ioReadall
			} else {
				ioReadall = ioutil.ReadAll
			}

			data, err := test.client.request(test.ctx, test.method, test.path, test.query, test.payload)
			assert.Equal(t, test.err, err)

			if err != nil {
				assert.Empty(t, data)
				return
			}

			var value map[string]interface{}
			err = json.Unmarshal([]byte(data), &value)
			assert.Nil(t, err)
			assert.Contains(t, value["headers"], "Authorization")
		})
	}
}

func TestListOptions_Check(t *testing.T) {
	tests := []struct {
		options ListOptions
		err     error
	}{
		{
			options: ListOptions{
				Type: "invalid",
			},
			err: errors.New("invalid type"),
		},
		{
			options: ListOptions{
				Type:  EmailTypeInbox,
				Order: "invalid",
			},
			err: errors.New("invalid order"),
		},
		{
			options: ListOptions{
				Type:  EmailTypeInbox,
				Order: OrderDesc,
			},
			err: nil,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := test.options.check()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestClientList(t *testing.T) {
	tests := []struct {
		client  Client
		options ListOptions
		args    map[string]interface{}
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			options: ListOptions{
				Type:  EmailTypeInbox,
				Order: OrderDesc,
			},
			args: map[string]interface{}{
				"type":  "inbox",
				"order": "desc",
			},
			err: nil,
		},
		{
			options: ListOptions{},
			err:     errors.New("invalid type"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.List(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)

			assert.Equal(t, test.args, values["args"])
		})
	}
}
