package duckdb

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
)

const driverName = "duckdb"

type Config struct {
	dataSource string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithDataSource(dataSource string) *Config {
	c.dataSource = dataSource
	return c
}

type Client struct {
	db *sql.DB
}

// New creates a new DuckDB client.
// The client is configured using the provided context.
func New(ctx context.Context, cfg *Config) (*Client, error) {
	conn, err := sql.Open(driverName, cfg.dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Client{
		db: conn,
	}, nil
}

// NewWithConn creates a new DuckDB client using an existing database connection.
// This is useful when you want to use a pre-established connection.
// The caller is responsible for closing the connection.
func NewWithConn(conn *sql.DB) *Client {
	return &Client{
		db: conn,
	}
}
