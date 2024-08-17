package meilisearch

import (
	"time"

	"github.com/meilisearch/meilisearch-go"
)

type Config struct {
	Host    *string
	ApiKey  *string
	Timeout *time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) WithHost(host string) *Config {
	c.Host = &host
	return c
}

func (c *Config) WithApiKey(apiKey string) *Config {
	c.ApiKey = &apiKey
	return c
}

func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = &timeout
	return c
}

type Client struct {
	client *meilisearch.Client
}

func New(config *Config) *Client {
	host := "http://localhost:7700"
	if config.Host != nil {
		host = *config.Host
	}
	apiKey := ""
	if config.ApiKey != nil {
		apiKey = *config.ApiKey
	}
	timeout := 5 * time.Second
	if config.Timeout != nil {
		timeout = *config.Timeout
	}

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:    host,
		APIKey:  apiKey,
		Timeout: timeout,
	})

	return &Client{client: client}
}
