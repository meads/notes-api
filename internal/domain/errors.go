package domain

import "errors"

type authError struct{ msg string }

func (e *authError) Error() string { return e.msg }

var (
	ErrInvalidCredentials = &authError{msg: "invalid username or password"}
	ErrTokenGeneration    = &authError{msg: "error token generation"}
	ErrSessionCreation    = &authError{msg: "error creating session"}
	ErrSessionRevoked     = &authError{msg: "error session revoked"}

	ErrInvalidToken = &authError{msg: "token is invalid"}
	// ErrInvalidSigningMethod = &authError{msg: "invalid signing method"}
	// ErrTokenExpired         = &authError{msg: "token is expired"}
)

var ErrUserNotFound = errors.New("user not found")
