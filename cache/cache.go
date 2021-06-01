package cache

import "errors"

var (
	ErrNotFound = errors.New("Not Found")
	ErrInternal = errors.New("Internal")
)
