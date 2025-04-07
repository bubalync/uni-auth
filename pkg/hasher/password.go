package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultCost = bcrypt.DefaultCost
)

type PasswordHasher interface {
	Hash(password string) ([]byte, error)
}

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptHasher(opts ...Option) *BcryptPasswordHasher {
	h := &BcryptPasswordHasher{
		cost: defaultCost,
	}

	// Custom options
	for _, opt := range opts {
		opt(h)
	}
	
	return h
}

func (h *BcryptPasswordHasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.cost)
}
