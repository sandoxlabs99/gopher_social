package store

import (
	"context"
	"database/sql"

	"github.com/sandoxlabs99/gopher_social/internal/models"
)

type RoleStore struct {
	db *sql.DB
}

func (s *RoleStore) GetByName(ctx context.Context, roleName string) (*models.Role, error) {
	query := `
		SELECT id, name, level, description FROM roles
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	r := &models.Role{}

	err := s.db.QueryRowContext(ctx, query, roleName).Scan(
		&r.ID,
		&r.Name,
		&r.Level,
		&r.Description,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}
