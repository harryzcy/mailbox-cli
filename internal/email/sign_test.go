package email

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestSignSDKRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	assert.Nil(t, err)

	err = SignSDKRequest(context.Background(), req, &SignSDKRequestOptions{
		Credentials: aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     "accessKeyID",
					SecretAccessKey: "secretAccessKey",
				}, nil
			},
		),
		Payload: []byte("payload"),
		Region:  "us-east-1",
	})
	assert.Nil(t, err)

	assert.NotEmpty(t, req.Header.Get("Authorization"))
	assert.IsType(t, "", req.Header.Get("Authorization"))
	assert.NotEmpty(t, req.Header.Get("X-Amz-Date"))
}

func TestSignSDKRequest_Error(t *testing.T) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	assert.Nil(t, err)

	err = SignSDKRequest(context.Background(), req, &SignSDKRequestOptions{
		Credentials: aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{}, errors.New("error")
			},
		),
	})
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("error"), err)

	err = SignSDKRequest(context.Background(), req, &SignSDKRequestOptions{})
	assert.Equal(t, ErrMissingCredentials, err)
}
