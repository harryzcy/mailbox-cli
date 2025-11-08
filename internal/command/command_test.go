package command

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/email"
	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T, handlerFunc http.HandlerFunc) *httptest.Server {
	ts := httptest.NewServer(handlerFunc)
	t.Cleanup(func() {
		ts.Close()
	})
	return ts
}

func TestGet(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

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
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := List(ListOptions{
		APIID:    "",
		Region:   "",
		Endpoint: ts.URL,
		Verbose:  false,
		Type:     "inbox",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestTrash(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Trash(TrashOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  ts.URL,
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestUntrash(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Untrash(UntrashOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  ts.URL,
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestDelete(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Delete(DeleteOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  ts.URL,
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestCreate(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Create(CreateOptions{
		APIID:        "",
		Region:       "",
		Endpoint:     ts.URL,
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
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestSave(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Save(SaveOptions{
		MessageID:    "messageID",
		APIID:        "",
		Region:       "",
		Endpoint:     ts.URL,
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
	assert.True(t, received, "Expected request to be received by the test server")
}

func TestSend(t *testing.T) {
	received := false
	ts := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Request received")
		assert.Nil(t, err)
		received = true
	})

	_, err := Send(SendOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  ts.URL,
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
	assert.True(t, received, "Expected request to be received by the test server")
}
