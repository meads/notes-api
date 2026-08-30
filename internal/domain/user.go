package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	CreatedAt time.Time
}

type ListUsersParams struct {
	Limit  int32
	Offset int32
}

type UpdateUserPasswordParams struct {
	ID              int64
	Username        string
	CurrentPassword string
	NewPassword     string
}

type UserRepository interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (*User, error)
	ListUsers(ctx context.Context, params ListUsersParams) ([]User, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	UserExists(ctx context.Context, id int64) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}
