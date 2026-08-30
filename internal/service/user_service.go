package service

import (
	"context"
	"fmt"

	domain "github.com/meads/notes-api/internal/domain"
	security "github.com/meads/notes-api/internal/security"
)

// implements handler.UserServicer
type UserService struct {
	userRepo domain.UserRepository
	hasher   security.Hasher
}

func NewUserService(repo domain.UserRepository, hasher security.Hasher) *UserService {
	return &UserService{
		userRepo: repo,
		hasher:   hasher,
	}
}

func (u *UserService) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.DeleteUser(ctx, id)
}

func (u *UserService) ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error) {
	return u.userRepo.ListUsers(ctx, params)
}
func (u *UserService) ChangePassword(ctx context.Context, params domain.UpdateUserPasswordParams) error {
	user, err := u.userRepo.GetUser(ctx, params.ID)
	if err != nil {
		return err
	}

	err = u.hasher.ComparePassword(user.Password, params.CurrentPassword)
	if err != nil {
		// current password invalid
		return domain.ErrInvalidCredentials
	}

	newPasswordHash, err := u.hasher.HashPassword(params.NewPassword)
	if err != nil {
		return fmt.Errorf("error creating password hash: %w", err)
	}

	err = u.userRepo.UpdateUserPassword(ctx, user.ID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	return nil
}
