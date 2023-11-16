package memory

import (
	"KeyGenerationService/internal/repository"
	"sync"
)

// InMemoryDB mocks the database for Key Generation Service.
type InMemoryDB struct {
	// Since our read and write are concurrent, use sync.Map instead of normal map and locks.
	Keys     sync.Map
	UsedKeys sync.Map
}

// New creates a new instance of InMemoryDB.
func New() (*InMemoryDB, error) {
	return &InMemoryDB{
		Keys:     sync.Map{},
		UsedKeys: sync.Map{},
	}, nil
}

// KeyExist checks whether a key exist within InMemoryDB.
func (i *InMemoryDB) KeyExist(key string) (bool, error) {
	if _, ok := i.Keys.Load(key); ok {
		return true, nil
	}
	if _, ok := i.UsedKeys.Load(key); ok {
		return true, nil
	}
	return false, repository.ErrKeyNotFound
}

// WriteKey stores the given key to InMemoryDB.
func (i *InMemoryDB) WriteKey(key string) error {
	i.Keys.Store(key, struct{}{})

	return nil
}

// GetKeys fetches an array of keys.
// The fetched keys are considered used and will be moved to UsedKeys for further usage.
func (i *InMemoryDB) GetKeys(requiredKeys int) ([]string, error) {
	// Cannot have negative or zero requiredKeys.
	if requiredKeys <= 0 {
		return []string{}, repository.ErrNegativeKey
	}

	// Cannot have requiredKeys greater than what we have in 'keys'.
	// This shouldn't happen since we always assume that we have enough keys in pool waiting.
	var mapLength int
	i.Keys.Range(func(_, _ any) bool {
		mapLength++
		return true
	})

	if requiredKeys > mapLength {
		return []string{}, repository.ErrKeyOutOfRange
	}

	// Create an array that stores all fetched keys.
	result := make([]string, requiredKeys)
	// Get keys randomly, and move used keys to used map.
	j := 0
	i.Keys.Range(func(key, value any) bool {
		if j == requiredKeys {
			return false
		}

		result[j] = key.(string)
		i.Keys.Delete(key)
		i.UsedKeys.Store(key, struct{}{})
		j++

		return true
	})
	return result, nil
}
