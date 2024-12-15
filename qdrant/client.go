package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
)

type Config struct {
	host        string
	port        int
	apiKey      *string
	useTls      *bool
	tlsConfig   *tls.Config
	grpcOptions []grpc.DialOption
}

func NewConfig(host string, port int) *Config {
	return &Config{
		host: host,
		port: port,
	}
}

func (c *Config) WithApiKey(apiKey string) *Config {
	c.apiKey = &apiKey
	return c
}

func (c *Config) WithTLS(tlsConfig *tls.Config) *Config {
	b := true
	c.useTls = &b
	c.tlsConfig = tlsConfig
	return c
}

func (c *Config) WithGrpcOptions(grpcOptions ...grpc.DialOption) *Config {
	c.grpcOptions = grpcOptions
	return c
}

type Client struct {
	client *qdrant.Client
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	qdrantCfg := &qdrant.Config{
		Host: cfg.host,
		Port: cfg.port,
	}

	if cfg.apiKey != nil {
		qdrantCfg.APIKey = *cfg.apiKey
	}
	if cfg.useTls != nil {
		qdrantCfg.UseTLS = *cfg.useTls
	}
	if cfg.tlsConfig != nil {
		qdrantCfg.TLSConfig = cfg.tlsConfig
	}
	if cfg.grpcOptions != nil {
		qdrantCfg.GrpcOptions = cfg.grpcOptions
	}

	client, err := qdrant.NewClient(qdrantCfg)
	if err != nil {
		return nil, fmt.Errorf("qdrant.NewClient: %w", err)
	}

	context.AfterFunc(ctx, func() {
		if err := client.Close(); err != nil {
			//fmt.Printf("Failed to close qdrant client: %v\n", err)
		}
	})

	return &Client{
		client: client,
	}, nil
}
