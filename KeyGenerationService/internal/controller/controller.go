package controller

import (
	"KeyGenerationService/internal/repository"
	"math/rand"
	"sync"
)

// letterBytes contains all possible characters for generating a Key.
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// KGS is the core for Key Generation Service.
type KGS struct {
	db repository.KGSDatabase
}

// NewKGS creates a new instance of KGS and generate keys concurrently to the database.
func NewKGS(size int, db repository.KGSDatabase) (*KGS, error) {
	kgs := &KGS{db: db}
	var wg sync.WaitGroup

	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				key := generateKey(4)
				exist, err := kgs.db.KeyExist(key)
				if err != nil {
					// log the error and continue
					continue
				}

				if !exist {
					_ = kgs.db.WriteKey(key)
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
	}()

	return kgs, nil
}

// generateKey generates the last four char in the shortenURL.
// Shortened URL should be in form: 'https://goShorten/1234'.
// Size of shortened URL is 16 bytes. Key in form '****' should have 62 ^ 4 variations.
// All variations take 16 * 62 ^ 4 bytes = 2,3642,1376 bytes, which is less than 0.3 GB.
// Implement a Key Generation Service using 1GB storage should be enough.
func generateKey(length int) string {
	res := make([]byte, length)

	for i := 0; i < length; i++ {
		res[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(res)
}

// GetKeys fetches an array of keys with length requiredKeys from the Key Generation Service database.
func (k *KGS) GetKeys(requiredKeys int) ([]string, error) {
	keys, err := k.db.GetKeys(requiredKeys)
	if err != nil {
		// Handler the error
		return nil, err
	}

	return keys, nil
}
