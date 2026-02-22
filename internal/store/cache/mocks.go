package cache

import (
	"context"

	"github.com/sandoxlabs99/gopher_social/internal/models"
	"github.com/stretchr/testify/mock"
)

func NewMockRedisStorage() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*models.User, error) {
	args := m.Called(userID)
	return nil, args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context, user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
