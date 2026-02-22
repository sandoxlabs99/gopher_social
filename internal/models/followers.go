package models

import "time"

type Follower struct {
	UserID     int64     `json:"userID"`
	FollowerID int64     `json:"followerID"`
	CreatedAt  time.Time `json:"createdAt"`
}
