package repository

import (
	"context"
	"errors"
)

// KGSDatabase is the interface that wraps writing and fetching keys from a Key Generation Service Database.
type KGSDatabase interface {
	KeyExist(context.Context, string) (bool, error)
	WriteKey(context.Context, string) error
	GetKeys(context.Context, int) ([]string, error)
}

var (
	ErrKeyNotFound   = errors.New("error desired key isn't found in database")
	ErrDatabaseError = errors.New("error malfunctioning of connecting to or using resource from a database")
	ErrKeyOOR        = errors.New("error key out of range")
)
