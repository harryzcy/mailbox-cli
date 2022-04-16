package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetErr(buf)

	c, err := rootCmd.ExecuteC()
	assert.Nil(t, err)
	assert.Equal(t, "CLI client for mailbox", c.Short)
	assert.Contains(t, buf.String(), "Usage:")
	assert.Contains(t, buf.String(), "Flags:")
}
