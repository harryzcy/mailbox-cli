package command

import "github.com/harryzcy/mailbox-cli/internal/email"

type GetOptions struct {
	// client options
	APIID    string
	Region   string
	Endpoint string
	Verbose  bool

	// request options
	MessageID string
}

func Get(options GetOptions) (string, error) {
	client := email.Client{
		APIID:    options.APIID,
		Region:   options.Region,
		Endpoint: options.Endpoint,
		Verbose:  options.Verbose,
	}

	result, err := client.Get(email.GetOptions{
		MessageID: options.MessageID,
	})

	return result, err
}

type ListOptions struct {
	// client options
	APIID    string
	Region   string
	Endpoint string
	Verbose  bool

	// request options
	Type       string
	Year       string
	Month      string
	Order      string // asc or desc (default)
	NextCursor string
}

func List(options ListOptions) (string, error) {
	client := email.Client{
		APIID:    options.APIID,
		Region:   options.Region,
		Endpoint: options.Endpoint,
		Verbose:  options.Verbose,
	}

	result, err := client.List(email.ListOptions{
		Type:       options.Type,
		Year:       options.Year,
		Month:      options.Month,
		Order:      options.Order,
		NextCursor: options.NextCursor,
	})

	return result, err
}

type TrashOptions struct {
	// client options
	APIID    string
	Region   string
	Endpoint string
	Verbose  bool

	// request options
	MessageID string
}

func Trash(options TrashOptions) (string, error) {
	client := email.Client{
		APIID:    options.APIID,
		Region:   options.Region,
		Endpoint: options.Endpoint,
		Verbose:  options.Verbose,
	}

	result, err := client.Trash(email.TrashOptions{
		MessageID: options.MessageID,
	})

	return result, err
}

type UntrashOptions struct {
	// client options
	APIID    string
	Region   string
	Endpoint string
	Verbose  bool

	// request options
	MessageID string
}

func Untrash(options UntrashOptions) (string, error) {
	client := email.Client{
		APIID:    options.APIID,
		Region:   options.Region,
		Endpoint: options.Endpoint,
		Verbose:  options.Verbose,
	}

	result, err := client.Untrash(email.UntrashOptions{
		MessageID: options.MessageID,
	})

	return result, err
}

type SendOptions struct {
	// client options
	APIID    string
	Region   string
	Endpoint string
	Verbose  bool

	// request options
	MessageID string
}

func Send(options SendOptions) (string, error) {
	client := email.Client{
		APIID:    options.APIID,
		Region:   options.Region,
		Endpoint: options.Endpoint,
		Verbose:  options.Verbose,
	}

	result, err := client.Send(email.SendOptions{
		MessageID: options.MessageID,
	})

	return result, err
}
