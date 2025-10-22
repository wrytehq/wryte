package models

import (
	"time"
)

type Session struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
