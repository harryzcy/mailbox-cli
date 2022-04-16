package email

import (
	"testing"

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
