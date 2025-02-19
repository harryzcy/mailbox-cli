package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestTrash(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"trash", "message-id"})

	commandTrash = func(_ command.TrashOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Trash an email", c.Short)
	assert.Equal(t, "result\n", buf.String())

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"trash"})
	_, err = rootCmd.ExecuteC()
	assert.NotNil(t, err)

	buf.Reset()
	commandTrash = func(_ command.TrashOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"trash", "message-id"})
	_, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
