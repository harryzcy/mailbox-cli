package email

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddQuery(t *testing.T) {
	tests := []struct {
		q         url.Values
		name      string
		value     string
		shouldAdd bool
	}{
		{
			q:         url.Values{},
			name:      "name",
			value:     "",
			shouldAdd: false,
		},
		{
			q:         url.Values{},
			name:      "name",
			value:     "value",
			shouldAdd: true,
		},
	}

	for _, test := range tests {
		addQuery(test.q, test.name, test.value)
		if test.shouldAdd {
			assert.NotEmpty(t, test.q)
			assert.Equal(t, test.value, test.q.Get(test.name))
		} else {
			assert.NotContains(t, test.q, test.name)
		}
	}
}
