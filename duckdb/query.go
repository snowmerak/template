package duckdb

import (
	"context"
	"fmt"
)

func (c *Client) CreatePerson(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, "CREATE TABLE person (name VARCHAR, age INTEGER)")
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
