package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "Handle mailbox APIs from the command line.", c.Short)
	assert.Contains(t, buf.String(), "Usage:")
	assert.Contains(t, buf.String(), "Flags:")
}

func TestExe(t *testing.T) {
	var exitCode int
	osExit = func(code int) { exitCode = code }

	rootCmd.SetArgs([]string{})
	Execute()
	assert.Equal(t, 0, exitCode)

	rootCmd.SetArgs([]string{"--help"})
	Execute()
	assert.Equal(t, 0, exitCode)

	rootCmd.SetArgs([]string{"invalid-command"})
	Execute()
	assert.Equal(t, 1, exitCode)
}
