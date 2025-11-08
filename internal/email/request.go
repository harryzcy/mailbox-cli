package email

import (
	"bytes"
	"encoding/json"
	"net/url"
)

func addQuery(q url.Values, name string, value string) {
	if value != "" {
		q.Add(name, value)
	}
}

func prettyResult(result []byte) (string, error) {
	var data map[string]any
	err := json.Unmarshal(result, &data)
	if err != nil {
		return "", err
	}

	// cannot use json.MarshalIndent because it escapes unicode characters
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
