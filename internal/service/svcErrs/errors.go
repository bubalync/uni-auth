package svcErrs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCannotParseToken   = errors.New("cannot parse token")
	ErrTokenIsExpired     = errors.New("token is expired")
	ErrCannotSignToken    = errors.New("cannot sign token")
	ErrAccessToCache      = errors.New("error access to cache service")

	ErrCannotCreateUser  = errors.New("cannot create user")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrCannotGetUser     = errors.New("cannot get user")
	ErrUserNotFound      = errors.New("user not found")
)
