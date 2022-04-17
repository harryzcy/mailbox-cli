package email

import (
	"encoding/json"
	"net/url"
)

func addQuery(q url.Values, name string, value string) {
	if value != "" {
		q.Add(name, value)
	}
}

func prettyResult(result []byte) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal(result, &data)
	if err != nil {
		return "", err
	}

	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
