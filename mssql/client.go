package mssql

import (
	"context"
	"database/sql/driver"
	"fmt"
	mssql "github.com/microsoft/go-mssqldb"
)

type Config struct {
	connectionString string
	connectionNumber int
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithConnectionString(connectionString string) *Config {
	c.connectionString = connectionString
	return c
}

func (c *Config) WithConnectionNumber(connectionNumber int) *Config {
	c.connectionNumber = connectionNumber
	return c
}

type Client struct {
	connector *mssql.Connector
}

func New(connectionString string) (*Client, error) {
	connector, err := mssql.NewConnector(connectionString)
	if err != nil {
		return nil, err
	}

	return &Client{
		connector: connector,
	}, nil
}

type Connection struct {
	conn driver.Conn
}

func (c *Client) Connect(ctx context.Context) (*Connection, error) {
	conn, err := c.connector.Connect(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Connection{
		conn: conn,
	}, nil
}
