package clickhouse

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Config struct {
	addr               []string
	database           *string
	username           *string
	password           *string
	insecureSkipVerify *bool

	clientInfo clickhouse.ClientInfo
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithAddr(addr ...string) *Config {
	c.addr = addr
	return c
}

func (c *Config) WithDatabase(database string) *Config {
	c.database = &database
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

func (c *Config) WithInsecureSkipVerify(insecureSkipVerify bool) *Config {
	c.insecureSkipVerify = &insecureSkipVerify
	return c
}

func (c *Config) AddClientInfo(name string, version string) *Config {
	c.clientInfo.Products = append(c.clientInfo.Products, struct {
		Name    string
		Version string
	}{Name: name, Version: version})
	return c
}

type Client struct {
	conn driver.Conn
}

func New(cfg Config) (*Client, error) {
	opt := clickhouse.Options{}
	if cfg.addr != nil {
		opt.Addr = cfg.addr
	}
	if cfg.database != nil {
		opt.Auth.Database = *cfg.database
	}
	if cfg.username != nil {
		opt.Auth.Username = *cfg.username
	}
	if cfg.password != nil {
		opt.Auth.Password = *cfg.password
	}
	if cfg.insecureSkipVerify != nil {
		opt.TLS.InsecureSkipVerify = *cfg.insecureSkipVerify
	}
	if len(cfg.clientInfo.Products) > 0 {
		opt.ClientInfo = cfg.clientInfo
	}

	conn, err := clickhouse.Open(&opt)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}

func NewWithConn(conn driver.Conn) *Client {
	return &Client{
		conn: conn,
	}
}
