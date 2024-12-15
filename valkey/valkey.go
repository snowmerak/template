package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

// Valkey is a Valkey client.
// It should be created using the New function.
type Valkey struct {
	conn valkey.Client
}

type Config struct {
	Addr     []string `json:"addr" yaml:"addr" xml:"addr" env:"addr"`
	Username string   `json:"username" yaml:"username" xml:"username" env:"username"`
	Password string   `json:"password" yaml:"password" xml:"password" env:"password"`
}

// New creates a new Valkey client.
// It returns an error if the client cannot be created.
// The caller is responsible for closing the client.
// The client is closed when the context is done.
func New(ctx context.Context, cfg *Config) (*Valkey, error) {
	conn, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: cfg.Addr,
		Username:    cfg.Username,
		Password:    cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("valkey.New: %w", err)
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Valkey{
		conn: conn,
	}, nil
}

// NewWithConn creates a new Valkey client with a custom connection.
// The caller is responsible for closing the client.
func NewWithConn(conn valkey.Client) *Valkey {
	return &Valkey{
		conn: conn,
	}
}

// Close closes the Valkey client.
// Close should be called when the client is no longer needed.
// But it is not necessary to call Close after a call to New.
// Because the context passed to New will close the connection pool when it is done.
func (r *Valkey) Close() {
	r.conn.Close()
}
