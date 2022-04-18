package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	_, err := Get(GetOptions{
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		MessageID: "messageID",
	})

	assert.Nil(t, err)
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
		APIID:    "",
		Region:   "",
		Endpoint: "https://httpbin.org/anything",
		Verbose:  false,
		Subject:  "subject",
		From:     []string{"from"},
		To:       []string{"to"},
		Cc:       []string{"cc"},
		Bcc:      []string{"bcc"},
		ReplyTo:  []string{"replyTo"},
		Body:     "body",
		Text:     "text",
		HTML:     "html",
	})

	assert.Nil(t, err)
}

func TestSave(t *testing.T) {
	_, err := Save(SaveOptions{
		MessageID: "messageID",
		APIID:     "",
		Region:    "",
		Endpoint:  "https://httpbin.org/anything",
		Verbose:   false,
		Subject:   "subject",
		From:      []string{"from"},
		To:        []string{"to"},
		Cc:        []string{"cc"},
		Bcc:       []string{"bcc"},
		ReplyTo:   []string{"replyTo"},
		Body:      "body",
		Text:      "text",
		HTML:      "html",
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
