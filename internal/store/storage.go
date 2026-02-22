package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sandoxlabs99/gopher_social/internal/models"
	"github.com/sandoxlabs99/gopher_social/internal/utils"
)

var (
	ErrNotFound          error = errors.New("resource not found")
	ErrUpdateConflict          = errors.New("update conflict")
	ErrDuplicateKey            = errors.New("resource already exists")
	QueryTimeoutDuration       = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *models.Post) error
		GetByID(context.Context, int64) (*models.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *models.Post) error
		GetUserFeed(context.Context, int64, utils.PaginatedFeedQuery) ([]models.PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *models.User) error
		GetByID(context.Context, int64) (*models.User, error)
		GetByEmail(context.Context, string) (*models.User, error)
		CreateAndInvite(context.Context, *models.User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		Create(context.Context, *models.Comment) error
		GetByPostID(context.Context, int64) ([]models.Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		UnFollow(ctx context.Context, unfollowedID, userID int64) error
	}
	Roles interface {
		GetByName(ctx context.Context, roleName string) (*models.Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
