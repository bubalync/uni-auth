package svcErrs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCannotParseToken   = errors.New("cannot parse token")
	ErrCannotSignToken    = errors.New("cannot sign token")
	ErrAccessToCache      = errors.New("access to cache service")

	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUserNotFound      = errors.New("user not found")
)
