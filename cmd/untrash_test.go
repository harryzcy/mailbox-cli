package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestUntrash(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"untrash", "message-id"})

	commandUntrash = func(options command.UntrashOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Untrash an email", c.Short)
	assert.Equal(t, "result\n", buf.String())

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"untrash"})
	c, err = rootCmd.ExecuteC()
	assert.NotNil(t, err)

	buf.Reset()
	commandUntrash = func(options command.UntrashOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"untrash", "message-id"})
	c, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
