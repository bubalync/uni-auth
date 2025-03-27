package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	PasswordHash     []byte    `json:"-"`
	IsActive         bool      `json:"is_active"`
	LastLoginAttempt time.Time `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
