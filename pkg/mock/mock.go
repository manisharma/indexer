package mock

import (
	"context"
	"indexer/pkg/models"
	"indexer/pkg/store"
)

// Mock/Dummy implementation of store.Repository to satisfy handler.HTTP
type Store struct{}

// Store implements store.Repository
var _ store.Repository = &Store{}

// return new mock implementation
func New() *Store {
	return &Store{}
}

func (s *Store) Create(ctx context.Context, e models.Epoch) error {
	return nil
}

func (s *Store) Get(ctx context.Context) ([]models.Epoch, error) {
	return []models.Epoch{}, nil
}

func (s *Store) KeepOnlyTop5(ctx context.Context, epochNumber uint64) error {
	return nil
}
