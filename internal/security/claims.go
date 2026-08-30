package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meads/notes-api/internal/domain"
)

type UserClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Usage    string `json:"usage"`
	*jwt.RegisteredClaims
}

func NewUserClaims(userID int64, username string, usage string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Error generating token id: %w", err)
	}
	return &UserClaims{
		UserID:   userID,
		Username: username,
		Usage:    usage,
		RegisteredClaims: &jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

type Tokener interface {
	GenerateToken(userID int64, username string, usage string, duration time.Duration) (string, *UserClaims, error)
	VerifyToken(tokenString string) (*UserClaims, error)
}

type TokenManager struct {
	secretKey string
}

func NewTokenManager(secretKey string) Tokener {
	return &TokenManager{
		secretKey: secretKey,
	}
}

func (c *TokenManager) GenerateToken(userID int64, username string, usage string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(userID, username, usage, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (c *TokenManager) VerifyToken(tokenStr string) (*UserClaims, error) {
	// Parse directly into the custom claims struct
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// invalid token signing method
			return nil, domain.ErrInvalidToken
		}
		return []byte(c.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrInvalidToken
		}
		return nil, domain.ErrInvalidToken
	}

	if token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
