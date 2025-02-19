package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harryzcy/mailbox-cli/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"save", "messageID"})

	commandSave = func(_ command.SaveOptions) (string, error) {
		return "result", nil
	}
	var exitCode int
	osExit = func(code int) { exitCode = code }

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Save a draft email", c.Short)
	assert.Equal(t, "result\n", buf.String())
	assert.Equal(t, 0, exitCode)

	// error
	buf.Reset()
	rootCmd.SetArgs([]string{"save"})
	_, err = rootCmd.ExecuteC()
	assert.NotNil(t, err)

	buf.Reset()
	commandSave = func(_ command.SaveOptions) (string, error) {
		return "result", errors.New("error")
	}
	rootCmd.SetArgs([]string{"save", "messageID"})
	_, err = rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, 1, exitCode)
	assert.Equal(t, "error\n", buf.String())
}
