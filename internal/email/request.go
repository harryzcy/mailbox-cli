package email

import "net/url"

func addQuery(q url.Values, name string, value string) {
	if value != "" {
		q.Add(name, value)
	}
}
