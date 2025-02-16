package adapter

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"porter"
)

type InMemoryRepository struct {
	data map[string]*porter.Port
	m    sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]*porter.Port),
	}
}

func (r *InMemoryRepository) GetPort(_ context.Context, key string) (*porter.Port, error) {
	r.m.Lock()
	defer r.m.Unlock()

	value, ok := r.data[key]
	if !ok {
		return nil, fmt.Errorf("port not found: %s", key)
	}

	return value, nil
}

func (r *InMemoryRepository) BatchCreateOrUpdatePort(ctx context.Context, data map[string]*porter.Port) error {
	r.m.Lock()
	defer r.m.Unlock()

	for key, value := range data {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		r.data[key] = value
	}

	log.Debug().Int("ports_count", len(r.data)).Msg("Current ports count.")

	return nil
}
