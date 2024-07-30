package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/snowmerak/template/postgres/queries"
)

// Client is a PostgreSQL client.
type Client struct {
	pool *pgxpool.Pool
}

// New creates a new client for the given connection string.
// The connection string should be a PostgreSQL connection string.
// New returns an error if it fails to create a new client.
// The client will automatically close the connection pool when the context is done.
func New(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	context.AfterFunc(ctx, func() {
		pool.Close()
	})

	return &Client{
		pool: pool,
	}, nil
}

// NewWithPool creates a new client with a custom connection pool.
// NewWithPool returns an error if the pool is nil.
// The caller is responsible for closing the pool when it is no longer needed.
func NewWithPool(pool *pgxpool.Pool) (*Client, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool is nil")
	}

	return &Client{
		pool: pool,
	}, nil
}

// Close closes the client's database connection pool.
// Close should be called when the client is no longer needed.
// But it is not necessary to call Close after a call to New.
// Because the context passed to New will close the connection pool when it is done.
func (c *Client) Close() {
	c.pool.Close()
}

// acquireConn acquires a connection from the pool.
// acquireConn returns a new Queries object that uses the acquired connection.
// acquireConn returns an error if it fails to acquire a connection.
func (c *Client) acquireConn(ctx context.Context) (*queries.Queries, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("pool.Acquire: %w", err)
	}

	return queries.New(conn), nil
}

// Write some methods here to interact with the database.
