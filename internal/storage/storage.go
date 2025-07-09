package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrLinkExists    = errors.New("link already exists")
	ErrLinkNotFound  = errors.New("link not found")
	ErrObjectExpired = errors.New("object expired")
	ErrCacheMiss     = errors.New("link not found in the cache")
)
