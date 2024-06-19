package redis

import (
	"context"
	"fmt"
	"github.com/redis/rueidis"
)

// Redis is a Redis client.
// It should be created using the New function.
type Redis struct {
	conn rueidis.Client
}

type Config struct {
	Addr     []string `json:"addr" yaml:"addr" xml:"addr" env:"addr"`
	Username string   `json:"username" yaml:"username" xml:"username" env:"username"`
	Password string   `json:"password" yaml:"password" xml:"password" env:"password"`
}

// New creates a new Redis client.
// It returns an error if the client cannot be created.
// The caller is responsible for closing the client.
// The client is closed when the context is done.
func New(ctx context.Context, cfg *Config) (*Redis, error) {
	conn, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: cfg.Addr,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("redis.New: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Redis{
		conn: conn,
	}, nil
}

// Close closes the Redis client.
// Close should be called when the client is no longer needed.
// But it is not necessary to call Close after a call to New.
// Because the context passed to New will close the connection pool when it is done.
func (r *Redis) Close() {
	r.conn.Close()
}
