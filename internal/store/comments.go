package store

import (
	"context"
	"database/sql"
	"gopher_social/internal/models"
)

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]models.Comment, error) {
	query := `
		SELECT 
			c.id, c.post_id, c.user_id, c.content, c.created_at,
			u.id, u.first_name, u.last_name, u.username, u.email, u.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}

	for rows.Next() {
		var c models.Comment
		c.User = models.User{}
		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt,
			&c.User.ID, &c.User.FirstName, &c.User.LastName,
			&c.User.Username, &c.User.Email, &c.User.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *models.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
