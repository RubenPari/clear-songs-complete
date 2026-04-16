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

func (m *MockAIRepository) ResolveArtistGenres(ctx context.Context, lookups []shared.AIArtistLookup) (map[string]string, error) {
	args := m.Called(ctx, lookups)
	var out map[string]string
	if v := args.Get(0); v != nil {
		out = v.(map[string]string)
	}
	return out, args.Error(1)
}

var _ shared.AIRepository = (*MockAIRepository)(nil)
