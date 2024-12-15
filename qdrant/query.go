package qdrant

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func (c *Client) CreateCollection(ctx context.Context, name string, size uint64, distance qdrant.Distance) (*qdrant.CollectionInfo, error) {
	if err := c.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     size,
			Distance: distance,
		}),
	}); err != nil {
		return nil, fmt.Errorf("CreateCollection: %w", err)
	}

	info, err := c.client.GetCollectionInfo(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("GetCollectionInfo: %w", err)
	}

	return info, nil
}

func (c *Client) DeleteCollection(ctx context.Context, name string) error {
	if err := c.client.DeleteCollection(ctx, name); err != nil {
		return fmt.Errorf("DeleteCollection: %w", err)
	}

	return nil
}

type Point struct {
	ID      uint64
	Vectors []float32
	Payload map[string]interface{}
}

func (c *Client) Upsert(ctx context.Context, collection string, data ...Point) (*qdrant.UpdateResult, error) {
	ups := make([]*qdrant.PointStruct, 0, len(data))
	for _, d := range data {
		ups = append(ups, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(d.ID),
			Vectors: qdrant.NewVectors(d.Vectors...),
			Payload: qdrant.NewValueMap(d.Payload),
		})
	}

	operationInfo, err := c.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
		Points:         ups,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert: %w", err)
	}

	return operationInfo, nil
}

type SearchOption struct {
}

type SearchResult struct {
	ID     uint64
	Score  float32
	Vector []float32
}

func (c *Client) Search(ctx context.Context, collection string, vector []float32, scoreThreshold float32, option *SearchOption) ([]*SearchResult, error) {
	searchResult, err := c.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection,
		Query:          qdrant.NewQuery(vector...),
		ScoreThreshold: &scoreThreshold,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	res := make([]*SearchResult, 0, len(searchResult))
	for _, r := range searchResult {
		rs := &SearchResult{}
		rs.ID = r.Id.GetNum()
		rs.Score = r.GetScore()
		rs.Vector = r.GetVectors().GetVector().GetData()
		// Parse payload here
		res = append(res, rs)
	}

	return res, nil
}

func (c *Client) Delete(ctx context.Context, collection string, ids ...uint64) (*qdrant.UpdateResult, error) {
	idsNum := make([]*qdrant.PointId, 0, len(ids))
	for _, id := range ids {
		idsNum = append(idsNum, qdrant.NewIDNum(id))
	}

	operationInfo, err := c.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: collection,
		Wait:           Box(true),
		Points:         qdrant.NewPointsSelector(idsNum...),
	})
	if err != nil {
		return nil, fmt.Errorf("delete: %w", err)
	}

	return operationInfo, nil
}
