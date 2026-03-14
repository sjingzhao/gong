package grpcx

// Client is a gRPC client
type Client struct {
}

// ClientOption is a client option
type ClientOption func(client *Client)

// NewClient creates a new gRPC client
func NewClient(opts ...ClientOption) *Client {
	cli := &Client{}

	for _, opt := range opts {
		opt(cli)
	}

	return cli
}
