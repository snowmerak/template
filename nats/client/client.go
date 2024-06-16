package client

import (
	"context"
	"github.com/nats-io/nats.go"
)

// Client is a wrapper around the nats.Conn struct
type Client struct {
	conn *nats.Conn
}

// New creates a new NATS client
// It returns an error if the client cannot be created
// The client is automatically closed when the context is done
func New(ctx context.Context) (*Client, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Client{conn: conn}, nil
}

// Publish publishes a message to a subject
func (c *Client) Publish(subject string, data []byte) error {
	return c.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject
func (c *Client) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

// Close closes the client
// It should be called when the client is no longer needed
// But it is automatically called when the context is done
func (c *Client) Close() {
	c.conn.Close()
}
