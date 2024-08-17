package meilisearch

import (
	"encoding/json"
	"fmt"
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

type SearchResult[T any] struct {
	Hits               []*T   `json:"hits"`
	Offset             int64  `json:"offset"`
	Limit              int64  `json:"limit"`
	EstimatedTotalHits int64  `json:"estimatedTotalHits"`
	ProcessingTimeMs   int64  `json:"processingTimeMs"`
	Query              string `json:"query"`
}

func Search[T any](c *Client, index string, query string) (*SearchResult[T], error) {
	res, err := c.client.Index(index).SearchRaw(query, &meilisearch.SearchRequest{
		Limit:            100,
		ShowRankingScore: true,
	})
	if err != nil {
		return nil, fmt.Errorf("meilisearch: searchRaw: %w", err)
	}

	sr := new(SearchResult[T])
	if err := json.Unmarshal(*res, sr); err != nil {
		return nil, fmt.Errorf("meilisearch: unmarshal: %w", err)
	}

	return sr, nil
}
