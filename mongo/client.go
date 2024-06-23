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
