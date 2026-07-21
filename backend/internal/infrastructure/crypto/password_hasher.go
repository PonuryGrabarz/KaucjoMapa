package crypto

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/alexedwards/argon2id"
)

type PasswordHasher struct{}

var _ auth.IPasswordHasher = (*PasswordHasher)(nil)

func NewHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (p *PasswordHasher) CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

func (p *PasswordHasher) HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}
