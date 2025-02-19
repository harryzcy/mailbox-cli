package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"create"})

	commandCreate = func(options command.CreateOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Create an email", c.Short)
	assert.Equal(t, "result\n", buf.String())
	assert.Equal(t, 0, exitCode)

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"create", "messageID"})
	_, err = rootCmd.ExecuteC()
	assert.NotNil(t, err)

	buf.Reset()
	commandCreate = func(_ command.CreateOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"create"})
	_, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
