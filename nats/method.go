package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
)

// Write your custom publisher and subscriber here

func (c *Client) Publish(ctx context.Context, subject string, data []byte) error {
	if err := c.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, subject string, handler func(data []byte)) error {
	sub, err := c.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}

	context.AfterFunc(ctx, func() {
		sub.Unsubscribe()
	})

	return nil
}
