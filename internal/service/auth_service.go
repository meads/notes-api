package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/meads/notes-api/internal/domain"
	security "github.com/meads/notes-api/internal/security"
)

// implements handler.AuthServicer
type AuthService struct {
	userRepo             domain.UserRepository
	sessionRepo          domain.SessionRepository
	tokener              security.Tokener
	hasher               security.Hasher
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	tokener security.Tokener,
	hasher security.Hasher,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		tokener:              tokener,
		hasher:               hasher,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (*domain.AuthSessionResult, error) {
	usernameExists, err := s.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("please choose another username")
	}

	password, err = s.hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.CreateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(user.ID, user.Username, "access", s.accessTokenDuration)
	if err != nil {
		return nil, domain.ErrTokenGeneration
	}

	refreshToken, refreshClaims, err := s.tokener.GenerateToken(user.ID, user.Username, "refresh", s.refreshTokenDuration)
	if err != nil {
		return nil, domain.ErrTokenGeneration
	}

	session, err := s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		return nil, domain.ErrSessionCreation
	}

	return &domain.AuthSessionResult{
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		SessionID:             session.ID,
		UserID:                user.ID,
	}, nil

}

func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.AuthSessionResult, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	err = s.hasher.ComparePassword(user.Password, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(user.ID, user.Username, "access", s.accessTokenDuration)
	if err != nil {
		return nil, domain.ErrTokenGeneration
	}

	refreshToken, refreshClaims, err := s.tokener.GenerateToken(user.ID, user.Username, "refresh", s.refreshTokenDuration)
	if err != nil {
		return nil, domain.ErrTokenGeneration
	}

	session, err := s.sessionRepo.CreateSession(ctx, domain.CreateSessionParams{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		return nil, domain.ErrSessionCreation
	}

	return &domain.AuthSessionResult{
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		SessionID:             session.ID,
		UserID:                user.ID,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	err := s.sessionRepo.DeleteSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("logout error, error deleting session %w", err)
	}
	return nil
}

func (s *AuthService) RenewAccessToken(ctx context.Context, refreshToken string) (*domain.RenewAccessTokenResult, error) {
	refreshClaims, err := s.tokener.VerifyToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	session, err := s.sessionRepo.GetSession(ctx, refreshClaims.RegisteredClaims.ID)
	if err != nil {
		return nil, err
	}

	if session.IsRevoked {
		return nil, domain.ErrSessionRevoked
	}

	accessToken, accessClaims, err := s.tokener.GenerateToken(
		refreshClaims.UserID, refreshClaims.Username, "access", s.accessTokenDuration)
	if err != nil {
		return nil, domain.ErrTokenGeneration
	}

	return &domain.RenewAccessTokenResult{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
	}, nil
}

func (s *AuthService) RevokeSession(ctx context.Context, sessionID string) error {
	err := s.sessionRepo.RevokeSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}
