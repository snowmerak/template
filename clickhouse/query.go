package clickhouse

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) FindExample(ctx context.Context, name string) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := c.conn.Query(ctx, "SELECT age FROM person WHERE name = ?", name)
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
