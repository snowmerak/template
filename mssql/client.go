package mssql

import (
	"context"
	"database/sql/driver"
	"errors"
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
	transact func(func(tx driver.Tx) error) error
	query    func(statement string, args []any, callback func(rows driver.Rows) error) error
	execute  func(query string, args []any, callback func(result driver.Result) error) error
}

func (c *Client) Connect(ctx context.Context) (*Connection, error) {
	conn, err := c.connector.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	transact := func(transact func(tx driver.Tx) error) (returnErr error) {
		tx, err := conn.Begin()
		if err != nil {
			returnErr = fmt.Errorf("failed to begin transaction: %w", err)
			return
		}

		commited := false
		defer func() {
			if !commited {
				if err := tx.Rollback(); err != nil {
					err = fmt.Errorf("failed to rollback transaction: %w", err)
					if returnErr != nil {
						returnErr = errors.Join(returnErr, err)
						return
					}
					returnErr = err
				}
			}
		}()

		if err := transact(tx); err != nil {
			returnErr = fmt.Errorf("failed to transact: %w", err)
			return
		}

		if err := tx.Commit(); err != nil {
			returnErr = fmt.Errorf("failed to commit transaction: %w", err)
			return
		}
		commited = true

		return nil
	}

	query := func(query string, args []any, callback func(rows driver.Rows) error) (returnErr error) {
		stmt, err := conn.Prepare(query)
		if err != nil {
			returnErr = fmt.Errorf("failed to prepare statement: %w", err)
			return
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				err = fmt.Errorf("failed to close statement: %w", err)

				if returnErr != nil {
					returnErr = errors.Join(returnErr, err)
					return
				}

				returnErr = err
				return
			}
		}()

		values := make([]driver.Value, len(args))
		for i, arg := range args {
			values[i] = arg
		}
		rows, err := stmt.Query(values)
		if err != nil {
			returnErr = fmt.Errorf("failed to query: %w", err)
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				err = fmt.Errorf("failed to close rows: %w", err)

				if returnErr != nil {
					returnErr = errors.Join(returnErr, err)
					return
				}

				returnErr = err
				return
			}
		}()

		if err := callback(rows); err != nil {
			return fmt.Errorf("failed to callback: %w", err)
		}

		return nil
	}

	execute := func(query string, args []any, callback func(result driver.Result) error) (returnErr error) {
		stmt, err := conn.Prepare(query)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer func() {
			if err := stmt.Close(); err != nil {
				err = fmt.Errorf("failed to close statement: %w", err)

				if returnErr != nil {
					returnErr = errors.Join(returnErr, err)
					return
				}

				returnErr = err
				return
			}
		}()

		values := make([]driver.Value, len(args))
		for i, arg := range args {
			values[i] = arg
		}
		result, err := stmt.Exec(values)
		if err != nil {
			returnErr = fmt.Errorf("failed to exec: %w", err)
			return
		}

		if err := callback(result); err != nil {
			returnErr = fmt.Errorf("failed to callback: %w", err)
			return
		}

		return nil
	}

	return &Connection{
		transact: transact,
		query:    query,
		execute:  execute,
	}, nil
}

func (c *Connection) Transact(transact func(tx driver.Tx) error) error {
	return c.transact(transact)
}

func (c *Connection) Query(statement string, args []any, callback func(rows driver.Rows) error) error {
	return c.query(statement, args, callback)
}

func (c *Connection) Execute(query string, args []any, callback func(result driver.Result) error) error {
	return c.execute(query, args, callback)
}
