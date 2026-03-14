package httpx

// Client is an http client
type Client struct {
}

// ClientOption is a client option
type ClientOption func(client *Client)

// NewClient creates a new http client
func NewClient(opts ...ClientOption) *Client {
	cli := &Client{}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}
