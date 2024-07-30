package nats

import (
	"context"
	"crypto/tls"
	"github.com/nats-io/nats.go"
	"strings"
)

type Config struct {
	url       *string
	username  *string
	password  *string
	token     *string
	tlsConfig *tls.Config
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithURL(url ...string) *Config {
	combinedURL := strings.Join(url, ",")
	c.url = &combinedURL
	return c
}

func (c *Config) WithUsername(username string) *Config {
	c.username = &username
	return c
}

func (c *Config) WithPassword(password string) *Config {
	c.password = &password
	return c
}

func (c *Config) WithToken(token string) *Config {
	c.token = &token
	return c
}

func (c *Config) WithTLSConfig(tlsConfig *tls.Config) *Config {
	c.tlsConfig = tlsConfig
	return c
}

// Client is a wrapper around the nats.Conn struct
type Client struct {
	conn *nats.Conn
}

// New creates a new NATS client
// It returns an error if the client cannot be created
// The client is automatically closed when the context is done
func New(ctx context.Context, cfg *Config) (*Client, error) {
	url := nats.DefaultURL
	if cfg.url != nil {
		url = *cfg.url
	}

	opt := nats.GetDefaultOptions()
	if cfg.username != nil {
		opt.User = *cfg.username
	}
	if cfg.password != nil {
		opt.Password = *cfg.password
	}
	if cfg.token != nil {
		opt.Token = *cfg.token
	}
	if cfg.tlsConfig != nil {
		opt.TLSConfig = cfg.tlsConfig
	}

	conn, err := nats.Connect(url, func(options *nats.Options) error {
		*options = opt
		return nil
	})
	if err != nil {
		return nil, err
	}

	cli := &Client{
		conn: conn,
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return cli, nil
}

// NewWithConn creates a new NATS client with a custom connection
// The caller is responsible for closing the connection
func NewWithConn(conn *nats.Conn) *Client {
	return &Client{
		conn: conn,
	}
}
