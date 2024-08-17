package meilisearch

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

func (c *Client) Query(index string, query string, options *meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	idx := c.client.Index(index)
	resp, err := idx.Search(query, options)
	if err != nil {
		return nil, fmt.Errorf("meilisearch: search failed: %w", err)
	}

	return resp, nil
}
