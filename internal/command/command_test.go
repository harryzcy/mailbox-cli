package command

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	received := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	}))
	defer ts.Close()

	_, err := Get(GetOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  ts.URL,
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestList(t *testing.T) {
	_, err := List(ListOptions{
		APIID:    "",
		Region:   "",
		Endpoint: "https://httpbin.org/anything",
		Verbose:  false,
		Type:     "inbox",
	})

	assert.Nil(t, err)
}

func TestTrash(t *testing.T) {
	_, err := Trash(TrashOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
}

func TestUntrash(t *testing.T) {
	_, err := Untrash(UntrashOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	_, err := Delete(DeleteOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	_, err := Create(CreateOptions{
		APIID:        "",
		Region:       "",
		Endpoint:     "https://httpbin.org/anything",
		Verbose:      false,
		Subject:      "subject",
		From:         []string{"from"},
		To:           []string{"to"},
		Cc:           []string{"cc"},
		Bcc:          []string{"bcc"},
		ReplyTo:      []string{"replyTo"},
		Text:         "text",
		HTML:         "html",
		GenerateText: email.GenerateTextAuto,
	})

	assert.Nil(t, err)
}

func TestSave(t *testing.T) {
	_, err := Save(SaveOptions{
		MessageID:    "messageID",
		APIID:        "",
		Region:       "",
		Endpoint:     "https://httpbin.org/anything",
		Verbose:      false,
		Subject:      "subject",
		From:         []string{"from"},
		To:           []string{"to"},
		Cc:           []string{"cc"},
		Bcc:          []string{"bcc"},
		ReplyTo:      []string{"replyTo"},
		Body:         "body",
		Text:         "text",
		HTML:         "html",
		GenerateText: email.GenerateTextAuto,
	})

	assert.Nil(t, err)
}

func TestSend(t *testing.T) {
	_, err := Send(SendOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
}
