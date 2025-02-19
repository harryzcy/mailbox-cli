package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"get", "message-id"})

	commandGet = func(_ command.GetOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Get an email by messageID", c.Short)
	assert.Equal(t, "result\n", buf.String())

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"get"})
	_, err = rootCmd.ExecuteC()
	assert.NotNil(t, err)

	buf.Reset()
	commandGet = func(_ command.GetOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"get", "message-id"})
	_, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
