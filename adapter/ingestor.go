package adapter

import (
	"context"
	"io"

	"porter/pkg/jsonstream"
)

type Ingestor[T any] struct{}

func (Ingestor[T]) Ingest(ctx context.Context, data io.Reader, doer func(key string, chunk *T) error) error {
	return jsonstream.Ingest[T](ctx, data, doer)
}

func NewIngestor[T any]() Ingestor[T] {
	return Ingestor[T]{}
}
