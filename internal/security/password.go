package security

import (
	"errors"

	"github.com/meads/notes-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword string, password string) error
}

type HashLib struct{}

func NewHasher() Hasher {
	return &HashLib{}
}

func (h *HashLib) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (h *HashLib) ComparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			return domain.ErrInvalidCredentials
		}
		return err
	}
	return nil
}
