package email

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func assertErrorsEqual(t *testing.T, actual, expected error) {
	if _, ok := actual.(interface{ Unwrap() []error }); ok {
		errs := actual.(interface{ Unwrap() []error }).Unwrap()
		assert.Contains(t, errs, expected)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func setupTestServer(t *testing.T, handlerFunc http.HandlerFunc) *httptest.Server {
	ts := httptest.NewServer(handlerFunc)
	t.Cleanup(func() {
		ts.Close()
	})
	return ts
}

func TestGetEndpoint(t *testing.T) {
	tests := []struct {
		client   Client
		endpoint string
	}{
		{
			client: Client{
				Endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com",
				Verbose:  true,
			},
			endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com",
		},
		{
			client: Client{
				APIID:   "api_id",
				Region:  "us-west-2",
				Verbose: true,
			},
			endpoint: "https://api_id.execute-api.us-west-2.amazonaws.com",
		},
	}

	for _, test := range tests {
		endpoint := test.client.getEndpoint()
		assert.Equal(t, test.endpoint, endpoint)
	}
}

func TestClient_Request(t *testing.T) {
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"headers": map[string]any{
				"Authorization": r.Header.Get("Authorization"),
			},
		}
		err := json.NewEncoder(w).Encode(response)
		assert.Nil(t, err)
	})

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Nanosecond) // used by one test case
	defer cancel()
	defer func() {
		ioReadall = io.ReadAll
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
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			method:    http.MethodGet,
			path:      "/get",
			query:     url.Values{},
			payload:   []byte(""),
			ioReadall: io.ReadAll,
			err:       nil,
		},
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			method:    http.MethodGet,
			path:      "/text",
			query:     url.Values{},
			payload:   []byte(""),
			ioReadall: io.ReadAll,
			err:       nil,
		},
		{
			client: Client{
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			err: errors.New("net/http: nil Context"),
		},
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, errors.New("error")
				}),
			},
			err: errors.New("error"),
		},
		{
			ctx: timeoutCtx,
			client: Client{
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			path:    "/get",
			query:   url.Values{},
			payload: []byte(""),
			err:     &url.Error{Op: "Get", URL: ts.URL + "/get", Err: context.DeadlineExceeded},
		},
		{
			ctx: context.Background(),
			client: Client{
				Endpoint: ts.URL,
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
			},
			path: "/get",
			ioReadall: func(_ io.Reader) ([]byte, error) {
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
				ioReadall = io.ReadAll
			}

			data, err := test.client.request(test.ctx, test.method, test.path, test.query, test.payload)
			assertErrorsEqual(t, err, test.err)
			if err != nil {
				assert.Empty(t, data)
				return
			}

			if test.path != "/text" {
				var value map[string]any
				err = json.Unmarshal([]byte(data), &value)
				assert.Nil(t, err)
				assert.Contains(t, value["headers"], "Authorization")
			}
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

func TestClient_List(t *testing.T) {
	defer func() {
		ioReadall = io.ReadAll
	}()

	tests := []struct {
		client    Client
		options   ListOptions
		ioReadall func(io.Reader) ([]byte, error)
		args      map[string]interface{}
		err       error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
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
			client: Client{
				Verbose: true,
			},
			options: ListOptions{},
			err:     errors.New("invalid type"),
		},
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: ListOptions{
				Type:  EmailTypeInbox,
				Order: OrderDesc,
			},
			ioReadall: func(_ io.Reader) ([]byte, error) {
				return nil, errors.New("error")
			},
			args: map[string]interface{}{
				"type":  "inbox",
				"order": "desc",
			},
			err: errors.New("error"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if test.ioReadall != nil {
				ioReadall = test.ioReadall
			} else {
				ioReadall = io.ReadAll
			}

			resp, err := test.client.List(test.options)
			assertErrorsEqual(t, err, test.err)
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

func TestGetOptions_Check(t *testing.T) {
	tests := []struct {
		options GetOptions
		err     error
	}{
		{
			options: GetOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: GetOptions{
				MessageID: "message-id",
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

func TestClient_Get(t *testing.T) {
	defer func() {
		ioReadall = io.ReadAll
	}()

	tests := []struct {
		client    Client
		options   GetOptions
		ioReadall func(io.Reader) ([]byte, error)
		err       error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: GetOptions{
				MessageID: "message-id",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: GetOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: GetOptions{
				MessageID: "message-id",
			},
			ioReadall: func(_ io.Reader) ([]byte, error) {
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
				ioReadall = io.ReadAll
			}

			resp, err := test.client.Get(test.options)
			assertErrorsEqual(t, err, test.err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestTrashOptions_Check(t *testing.T) {
	tests := []struct {
		options TrashOptions
		err     error
	}{
		{
			options: TrashOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: TrashOptions{
				MessageID: "message-id",
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

func TestClient_Trash(t *testing.T) {
	tests := []struct {
		client  Client
		options TrashOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: TrashOptions{
				MessageID: "message-id",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: TrashOptions{},
			err:     errors.New("invalid message id"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Trash(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestUntrashOptions_Check(t *testing.T) {
	tests := []struct {
		options UntrashOptions
		err     error
	}{
		{
			options: UntrashOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: UntrashOptions{
				MessageID: "message-id",
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

func TestClient_Untrash(t *testing.T) {
	tests := []struct {
		client  Client
		options UntrashOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: UntrashOptions{
				MessageID: "message-id",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: UntrashOptions{},
			err:     errors.New("invalid message id"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Untrash(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestDeleteOptions_Check(t *testing.T) {
	tests := []struct {
		options DeleteOptions
		err     error
	}{
		{
			options: DeleteOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: DeleteOptions{
				MessageID: "message-id",
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

func TestClient_Delete(t *testing.T) {
	tests := []struct {
		client  Client
		options DeleteOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: DeleteOptions{
				MessageID: "message-id",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: DeleteOptions{},
			err:     errors.New("invalid message id"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Delete(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestCreateOptions_Check(t *testing.T) {
	tests := []struct {
		options CreateOptions
		err     error
	}{
		{
			options: CreateOptions{
				Subject: "subject",
				From:    []string{"from"},
				To:      []string{"to"},
			},
			err: errors.New("invalid generate-text"),
		},
		{
			options: CreateOptions{
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "invalid",
			},
			err: errors.New("invalid generate-text"),
		},
		{
			options: CreateOptions{
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "on",
			},
			err: nil,
		},
		{
			options: CreateOptions{
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "off",
			},
			err: nil,
		},
		{
			options: CreateOptions{
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: GenerateTextAuto,
			},
			err: nil,
		},
		{
			options: CreateOptions{
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				Cc:           []string{"cc"},
				Bcc:          []string{"bcc"},
				ReplyTo:      []string{"reply-to"},
				Text:         "text",
				HTML:         "html",
				GenerateText: GenerateTextAuto,
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

func TestCreateOptions_LoadFile(t *testing.T) {
	tests := []struct {
		options CreateOptions
		err     error
	}{
		{
			options: CreateOptions{},
			err:     nil,
		},
		{
			options: CreateOptions{
				File: "../../test/data/email.json",
			},
			err: nil,
		},
		{
			options: CreateOptions{
				File: "invalid.json",
			},
			err: &os.PathError{Op: "open", Path: "invalid.json", Err: syscall.ENOENT},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := test.options.loadFile()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestClient_Create(t *testing.T) {
	tests := []struct {
		client  Client
		options CreateOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: CreateOptions{
				File: "../../test/data/email.json",
			},
			err: nil,
		},
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Verbose:  true,
			},
			options: CreateOptions{},
			err:     nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: CreateOptions{
				File: "invalid.json",
			},
			err: &os.PathError{Op: "open", Path: "invalid.json", Err: syscall.ENOENT},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Create(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestSaveOptions_Check(t *testing.T) {
	tests := []struct {
		options SaveOptions
		err     error
	}{
		{
			options: SaveOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: SaveOptions{
				MessageID: "message-id",
			},
			err: errors.New("invalid subject"),
		},
		{
			options: SaveOptions{
				MessageID: "message-id",
				Subject:   "subject",
			},
			err: errors.New("invalid from"),
		},
		{
			options: SaveOptions{
				MessageID: "message-id",
				Subject:   "subject",
				From:      []string{"from"},
			},
			err: errors.New("invalid to"),
		},
		{
			options: SaveOptions{
				MessageID: "message-id",
				Subject:   "subject",
				From:      []string{"from"},
				To:        []string{"to"},
			},
			err: errors.New("invalid generate-text"),
		},
		{
			options: SaveOptions{
				MessageID:    "message-id",
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "invalid",
			},
			err: errors.New("invalid generate-text"),
		},
		{
			options: SaveOptions{
				MessageID:    "message-id",
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "on",
			},
			err: nil,
		},
		{
			options: SaveOptions{
				MessageID:    "message-id",
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: "off",
			},
			err: nil,
		},
		{
			options: SaveOptions{
				MessageID:    "message-id",
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				GenerateText: GenerateTextAuto,
			},
			err: nil,
		},
		{
			options: SaveOptions{
				MessageID:    "message-id",
				Subject:      "subject",
				From:         []string{"from"},
				To:           []string{"to"},
				Cc:           []string{"cc"},
				Bcc:          []string{"bcc"},
				ReplyTo:      []string{"reply-to"},
				Body:         "body",
				Text:         "text",
				HTML:         "html",
				GenerateText: GenerateTextAuto,
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

func TestSaveOptions_LoadFile(t *testing.T) {
	tests := []struct {
		options SaveOptions
		err     error
	}{
		{
			options: SaveOptions{},
			err:     nil,
		},
		{
			options: SaveOptions{
				File: "../../test/data/email.json",
			},
			err: nil,
		},
		{
			options: SaveOptions{
				File: "invalid.json",
			},
			err: &os.PathError{Op: "open", Path: "invalid.json", Err: syscall.ENOENT},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := test.options.loadFile()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestClient_Save(t *testing.T) {
	tests := []struct {
		client  Client
		options SaveOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: SaveOptions{
				MessageID: "message-id",
				File:      "../../test/data/email.json",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: SaveOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			client: Client{
				Verbose: true,
			},
			options: SaveOptions{
				MessageID: "message-id",
				File:      "invalid.json",
			},
			err: &os.PathError{Op: "open", Path: "invalid.json", Err: syscall.ENOENT},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Save(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}

func TestSendOptions_Check(t *testing.T) {
	tests := []struct {
		options SendOptions
		err     error
	}{
		{
			options: SendOptions{},
			err:     errors.New("invalid message id"),
		},
		{
			options: SendOptions{
				MessageID: "message-id",
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

func TestClient_Send(t *testing.T) {
	tests := []struct {
		client  Client
		options SendOptions
		err     error
	}{
		{
			client: Client{
				Endpoint: "https://httpbin.org/anything",
				Credentials: aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
					return aws.Credentials{}, nil
				}),
				Verbose: true,
			},
			options: SendOptions{
				MessageID: "message-id",
			},
			err: nil,
		},
		{
			client: Client{
				Verbose: true,
			},
			options: SendOptions{},
			err:     errors.New("invalid message id"),
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			resp, err := test.client.Send(test.options)
			assert.Equal(t, test.err, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, resp)

			var values map[string]interface{}
			err = json.Unmarshal([]byte(resp), &values)
			assert.Nil(t, err)
		})
	}
}
