package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *client.Client
}

func Default(ctx context.Context) (*Client, error) {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		c.Close()
	})

	return &Client{client: c}, nil
}

type Config struct {
	endpoint   string
	httpClient *http.Client

	timeout time.Duration

	cacertPath string
	certPath   string
	keyPath    string

	userAgent string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithEndpoint(endpoint string) *Config {
	c.endpoint = endpoint
	return c
}

func (c *Config) WithHTTPClient(httpClient *http.Client) *Config {
	c.httpClient = httpClient
	return c
}

func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.timeout = timeout
	return c
}

func (c *Config) WithTlsCfg(cacertPath, certPath, keyPath string) *Config {
	c.cacertPath = cacertPath
	c.certPath = certPath
	c.keyPath = keyPath
	return c
}

func (c *Config) WithUserAgent(userAgent string) *Config {
	c.userAgent = userAgent
	return c
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}

	if cfg.endpoint != "" {
		opts = append(opts, client.WithHost(cfg.endpoint))
	}
	if cfg.httpClient != nil {
		opts = append(opts, client.WithHTTPClient(cfg.httpClient))
	}
	if cfg.timeout != 0 {
		opts = append(opts, client.WithTimeout(cfg.timeout))
	}
	if cfg.cacertPath != "" {
		opts = append(opts, client.WithTLSClientConfig(cfg.certPath, cfg.keyPath, cfg.cacertPath))
	}
	if cfg.userAgent != "" {
		opts = append(opts, client.WithUserAgent(cfg.userAgent))
	}

	c, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		c.Close()
	})

	return &Client{
		client: c,
	}, nil
}

func AuthString(username string, password string) string {
	builder := bytes.NewBuffer(nil)
	builder.WriteString("{\"username\":\"")
	builder.WriteString(username)
	builder.WriteString("\",\"password\":\"")
	builder.WriteString(password)
	builder.WriteString("\"}")

	return base64.URLEncoding.EncodeToString(builder.Bytes())
}

func (c *Client) Pull(ctx context.Context, imagePath string, authString string, writer io.Writer) error {
	reader, err := c.client.ImagePull(ctx, imagePath, image.PullOptions{
		RegistryAuth: authString,
	})
	if err != nil {
		return fmt.Errorf("pull image %s failed: %w", imagePath, err)
	}
	defer reader.Close()

loop:
	for {
		_, err := io.Copy(writer, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break loop
			}
			return fmt.Errorf("pull image %s failed: %w", imagePath, err)
		}
	}

	return nil
}

func (c *Client) Push(ctx context.Context, imagePath string, authString string, writer io.Writer) error {
	reader, err := c.client.ImagePush(ctx, imagePath, image.PushOptions{
		RegistryAuth: authString,
	})
	if err != nil {
		return fmt.Errorf("push image %s failed: %w", imagePath, err)
	}
	defer reader.Close()

loop:
	for {
		_, err := io.Copy(writer, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break loop
			}
			return fmt.Errorf("push image %s failed: %w", imagePath, err)
		}
	}

	return nil
}
