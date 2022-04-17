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
