package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list"})

	commandList = func(_ command.ListOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "List emails", c.Short)
	assert.Equal(t, "result\n", buf.String())

	// error
	buf.Reset()
	commandList = func(_ command.ListOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"list"})
	_, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
