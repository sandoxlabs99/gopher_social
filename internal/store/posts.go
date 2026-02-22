package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sandoxlabs99/gopher_social/internal/models"
	"github.com/sandoxlabs99/gopher_social/internal/utils"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *models.Post) error {
	query := `
	INSERT INTO posts (title, content, tags, user_id)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// use pg_sleep(10) or generate_series to simulate long running/sslow queries to test the context timeout

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, postID int64) (*models.Post, error) {
	var post models.Post

	query := `
		SELECT 
			id, title, content, tags, user_id, 
			created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&post.ID, &post.Title, &post.Content, pq.Array(&post.Tags),
		&post.UserID, &post.CreatedAt, &post.UpdatedAt, &post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// Performs Optimistic Locking/Concurrency
func (s *PostStore) Update(ctx context.Context, post *models.Post) error {
	query := `
		UPDATE posts
		SET 
			title = $1, 
			content = $2, 
			tags = $3,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $4 AND version = $5
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx, query, post.Title, post.Content,
		pq.Array(post.Tags), post.ID, post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrUpdateConflict
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq utils.PaginatedFeedQuery) ([]models.PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.title, p.content, p.tags,
			p.user_id, p.created_at, p.version, u.username,
			COALESCE(c.comment_count, 0) AS comments_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN followers f ON f.user_id = u.id AND f.follower_id = $1
		LEFT JOIN (
			SELECT post_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY post_id
		) c ON c.post_id = p.id
		WHERE (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND (p.tags @> $5 OR $5 = '{}')
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []models.PostWithMetadata
	for rows.Next() {
		var p models.PostWithMetadata
		err := rows.Scan(
			&p.ID, &p.Title, &p.Content,
			pq.Array(&p.Tags), &p.UserID, &p.CreatedAt,
			&p.Version, &p.User.Username, &p.CommentCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}
