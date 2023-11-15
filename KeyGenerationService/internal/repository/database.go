package repository

import "errors"

// KGSDatabase is the interface that wraps writing and fetching keys from a Key Generation Service Database.
type KGSDatabase interface {
	KeyExist(string) (bool, error)
	WriteKey(string) error
	GetKeys(int) ([]string, error)
}

var (
	ErrKeyNotFound = errors.New("error desired key isn't found in database")
	ErrNegativeKey = errors.New("error cannot have 0 or negative requiredKeys")
)
