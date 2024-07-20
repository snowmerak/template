package duckdb

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
)

const driverName = "duckdb"

type Client struct {
	db *sql.DB
}

func New(ctx context.Context) (*Client, error) {
	conn, err := sql.Open(driverName, "")
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
