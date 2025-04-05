package svcErrs

import "errors"

var (
	ErrInternal          = errors.New("internal server error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
