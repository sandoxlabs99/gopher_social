package store

import (
	"context"
	"database/sql"
	"gopher_social/internal/models"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(context.Context, *sql.Tx, *models.User) error {
	return nil
}

func (m *MockUserStore) GetByID(context.Context, int64) (*models.User, error) {
	return &models.User{
		ID: 1,
	}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (m *MockUserStore) CreateAndInvite(context.Context, *models.User, string, time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(context.Context, string) error {
	return nil
}

func (m *MockUserStore) Delete(context.Context, int64) error {
	return nil
}
