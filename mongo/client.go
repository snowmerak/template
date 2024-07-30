package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

type OptionApply func(*options.ClientOptions)

// New creates a new MongoDB client
// It returns an error if the client cannot be created
// The client is automatically closed when the context is done
func New(ctx context.Context, apply OptionApply) (*Client, error) {
	opt := options.Client()
	apply(opt)

	cli, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		cli.Disconnect(ctx)
	})

	return &Client{
		client: cli,
	}, nil
}

// NewWithConn creates a new MongoDB client with a custom connection
// The caller is responsible for closing the client
func NewWithConn(conn *mongo.Client) *Client {
	return &Client{
		client: conn,
	}
}
