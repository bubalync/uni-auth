package services

import "errors"

var (
	ErrInternal     = errors.New("internal server error")
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
