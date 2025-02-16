package porter

import (
	"context"
	"io"
)

//go:generate moq -rm -out mock_service.go . Service
type Service interface {
	// PortsFromStream creates or updates a port from a stream of data.
	PortsFromStream(ctx context.Context, data io.Reader) error

	// GetPort retrieves a port by its UN/LOCODE.
	GetPort(ctx context.Context, unloc string) (*Port, error)
}
