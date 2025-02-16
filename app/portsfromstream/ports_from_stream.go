package portsfromstream

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"porter"
)

//go:generate moq -rm -out mock_repository.go . Repository
type Repository interface {
	BatchCreateOrUpdatePort(ctx context.Context, data map[string]*porter.Port) error
}

//go:generate moq -rm -out mock_ingestor.go . Ingestor
type Ingestor interface {
	Ingest(ctx context.Context, data io.Reader, doer func(key string, port *porter.Port) error) error
}

func NewService(repo Repository, ingestor Ingestor, batchSize int) *PortsFromStreamService {
	return &PortsFromStreamService{
		repo:      repo,
		ingestor:  ingestor,
		batchSize: batchSize,
	}
}

type PortsFromStreamService struct {
	repo      Repository
	ingestor  Ingestor
	batchSize int // Adjust depending on repository and available memory.
}

func (s *PortsFromStreamService) PortsFromStream(ctx context.Context, data io.Reader) error {
	var counter int

	buffer := make(map[string]*porter.Port, s.batchSize)

	if err := s.ingestor.Ingest(ctx, data, func(key string, port *porter.Port) error {
		if err := port.Validate(); err != nil {
			log.Err(err).Interface("port", port).Msg("Failed to validate port.")

			return nil
		}

		counter++

		buffer[key] = port

		if counter%s.batchSize == 0 {
			if err := s.repo.BatchCreateOrUpdatePort(ctx, buffer); err != nil {
				return fmt.Errorf("failed to save batch: %w", err)
			}

			log.Info().Int("batch_size", s.batchSize).Msg("Saved batch.")

			counter = 0
			clear(buffer)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to ingest data: %w", err)
	}

	if counter > 0 {
		if err := s.repo.BatchCreateOrUpdatePort(ctx, buffer); err != nil {
			return fmt.Errorf("failed to save remaining ports: %w", err)
		}

		log.Info().Int("batch_size", counter).Msg("Saved remaining ports.")
	}

	return nil
}
