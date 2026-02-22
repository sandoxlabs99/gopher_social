package models

import (
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postID"`
	UserID    int64     `json:"userID"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}
