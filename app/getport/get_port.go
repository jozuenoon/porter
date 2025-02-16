package getport

import (
	"context"

	"porter"
)

//go:generate moq -rm -out mock_repository.go . Repository
type Repository interface {
	GetPort(ctx context.Context, unloc string) (*porter.Port, error)
}

func NewService(repo Repository) *GetPortService {
	return &GetPortService{
		repo: repo,
	}
}

type GetPortService struct {
	repo Repository
}

func (s *GetPortService) GetPort(ctx context.Context, unloc string) (*porter.Port, error) {
	return s.repo.GetPort(ctx, unloc)
}
