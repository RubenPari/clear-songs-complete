package mocks

import (
	"context"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/stretchr/testify/mock"
)

// MockAIRepository is a mock implementation of AIRepository
type MockAIRepository struct {
	mock.Mock
}

func (m *MockAIRepository) ResolveArtistGenre(ctx context.Context, artistName string) (string, error) {
	args := m.Called(ctx, artistName)
	return args.String(0), args.Error(1)
}

var _ shared.AIRepository = (*MockAIRepository)(nil)
