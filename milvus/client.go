package milvus

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	client client.Client
}

type Config struct {
	address        string
	username       string
	password       string
	dbName         string
	identifier     string
	enableTLSAuth  bool
	apiKey         string
	serverVersion  string
	dialOptions    []grpc.DialOption
	retryRateLimit *client.RetryRateLimitOption
	disableConn    bool
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) WithAddress(address string) *Config {
	cfg.address = address
	return cfg
}

func (cfg *Config) WithUsername(username string) *Config {
	cfg.username = username
	return cfg
}

func (cfg *Config) WithPassword(password string) *Config {
	cfg.password = password
	return cfg
}

func (cfg *Config) WithDBName(dbName string) *Config {
	cfg.dbName = dbName
	return cfg
}

func (cfg *Config) WithIdentifier(identifier string) *Config {
	cfg.identifier = identifier
	return cfg
}

func (cfg *Config) WithEnableTLSAuth(enableTLSAuth bool) *Config {
	cfg.enableTLSAuth = enableTLSAuth
	return cfg
}

func (cfg *Config) WithAPIKey(apiKey string) *Config {
	cfg.apiKey = apiKey
	return cfg
}

func (cfg *Config) WithServerVersion(serverVersion string) *Config {
	cfg.serverVersion = serverVersion
	return cfg
}

func (cfg *Config) WithDialOptions(dialOptions []grpc.DialOption) *Config {
	cfg.dialOptions = dialOptions
	return cfg
}

func (cfg *Config) WithRetryRateLimit(maxRetryAttempt uint, maxRetryBackoff time.Duration) *Config {
	cfg.retryRateLimit = &client.RetryRateLimitOption{
		MaxRetry:   maxRetryAttempt,
		MaxBackoff: maxRetryBackoff,
	}
	return cfg
}

func (cfg *Config) WithDisableConn(disableConn bool) *Config {
	cfg.disableConn = disableConn
	return cfg
}

func New(ctx context.Context, cfg *Config) (*Client, error) {
	cli, err := client.NewClient(ctx, client.Config{
		Address:        cfg.address,
		Username:       cfg.username,
		Password:       cfg.password,
		DBName:         cfg.dbName,
		Identifier:     cfg.identifier,
		EnableTLSAuth:  cfg.enableTLSAuth,
		APIKey:         cfg.apiKey,
		ServerVersion:  cfg.serverVersion,
		DialOptions:    cfg.dialOptions,
		RetryRateLimit: cfg.retryRateLimit,
		DisableConn:    cfg.disableConn,
	})
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	context.AfterFunc(ctx, func() {
		cli.Close()
	})

	return &Client{
		client: cli,
	}, nil
}

// NewWithConn creates a new client with an existing connection.
// The caller is responsible for closing the connection.
func NewWithConn(conn client.Client) *Client {
	return &Client{
		client: conn,
	}
}

func (c *Client) Close() {
	c.client.Close()
}
