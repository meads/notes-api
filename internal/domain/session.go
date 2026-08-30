package domain

import (
	"context"
	"time"
)

type Session struct {
	ID           string
	UserID       int64
	RefreshToken string
	IsRevoked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type CreateSessionParams struct {
	ID           string
	UserID       int64
	RefreshToken string
	IsRevoked    bool
	ExpiresAt    time.Time
}

type SessionRepository interface {
	CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
	RevokeSession(ctx context.Context, id string) error
}
