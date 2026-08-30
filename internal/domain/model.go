package domain

import "time"

type AuthSessionResult struct {
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	SessionID             string
	UserID                int64
}

type RenewAccessTokenResult struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}
