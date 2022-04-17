package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"send", "message-id"})

	commandSend = func(options command.SendOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Send an email", c.Short)
	assert.Equal(t, "result\n", buf.String())

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"send"})
	c, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Contains(t, buf.String(), "Please specify a messageID")

	buf.Reset()
	commandSend = func(options command.SendOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"send", "message-id"})
	c, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
