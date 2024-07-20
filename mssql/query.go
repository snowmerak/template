package mssql

import (
	"context"
	"database/sql/driver"
	"fmt"
)

func Values(args ...interface{}) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		values[i] = arg
	}
	return values
}

func (c *Client) Find(ctx context.Context, name string) ([]int64, error) {
	stmt, err := c.conn.PrepareContext(ctx, "SELECT age FROM person WHERE name = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	ages := make([]int64, 0, 4)
	for rows.Next() {
		var age int64
		if err := rows.Scan(&age); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		ages = append(ages, age)
	}

	return ages, nil
}
