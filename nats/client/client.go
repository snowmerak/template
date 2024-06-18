package client

import (
	"context"
	"github.com/nats-io/nats.go"
)

// Client is a wrapper around the nats.Conn struct
type Client[T any] struct {
	conn    *nats.Conn
	encoder func(T) ([]byte, error)
	decoder func([]byte) (T, error)
}

// New creates a new NATS client
// It returns an error if the client cannot be created
// The client is automatically closed when the context is done
func New[T any](ctx context.Context, encoder func(T) ([]byte, error), decoder func([]byte) (T, error)) (*Client[T], error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Client[T]{
		conn:    conn,
		encoder: encoder,
		decoder: decoder,
	}, nil
}

// Publish publishes a message to a subject
func (c *Client[T]) Publish(subject string, value T) error {
	data, err := c.encoder(value)
	if err != nil {
		return err
	}

	return c.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject
func (c *Client[T]) Subscribe(subject string, handler func(subject string, value T)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, func(msg *nats.Msg) {
		value, err := c.decoder(msg.Data)
		if err != nil {
			return
		}

		handler(msg.Subject, value)
	})
}

// Close closes the client
// It should be called when the client is no longer needed
// But it is automatically called when the context is done
func (c *Client[T]) Close() {
	c.conn.Close()
}
