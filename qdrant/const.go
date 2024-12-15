package qdrant

import "github.com/qdrant/go-client/qdrant"

func Box[T any](value T) *T {
	return &value
}

const (
	DistanceCosine    = qdrant.Distance_Cosine
	DistanceEuclidean = qdrant.Distance_Euclid
	DistanceDot       = qdrant.Distance_Dot
	DistanceManhattan = qdrant.Distance_Manhattan
)
