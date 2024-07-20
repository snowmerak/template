package mssql

import (
	"database/sql"
	"fmt"
	"github.com/microsoft/go-mssqldb/azuread"
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
	conn *sql.DB
}

type DriverName string

const (
	DriverNameMSSQL   DriverName = "mssql"
	DriverNameAzureAD DriverName = azuread.DriverName
)

func New(driver DriverName, connectionString string) (*Client, error) {
	conn, err := sql.Open(string(driver), connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	return &Client{
		conn: conn,
	}, nil
}
