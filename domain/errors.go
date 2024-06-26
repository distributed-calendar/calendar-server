package domain

import "errors"

var (
	ErrDataNotFound = errors.New("data not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrExternalUnavailable = errors.New("external resource is unavailable")
)
