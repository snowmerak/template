package mssql

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
)

func Values(args ...interface{}) []driver.Value {
	values := make([]driver.Value, len(args))
	for i, arg := range args {
		values[i] = arg
	}
	return values
}

func (c *Connection) Find(ctx context.Context, name string) ([]int64, error) {
	stmt, err := c.conn.Prepare("SELECT age FROM person WHERE name = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	rows, err := stmt.Query(Values(name))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	ages := make([]int64, 0, 4)
	for {
		var age int64
		if err := rows.Next(Values(&age)); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to read row: %w", err)
		}

		ages = append(ages, age)
	}

	return ages, nil
}
