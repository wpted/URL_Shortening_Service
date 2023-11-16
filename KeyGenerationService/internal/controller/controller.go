package controller

import (
	"KeyGenerationService/internal/repository"
	"context"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

// LetterBytes contains all possible characters for generating a Key.
const LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	ErrRepoError        = errors.New("error repo failed")
	ErrInvalidKeyLength = errors.New("error cannot have key length equal or smaller than 0")
	ErrGetKeysError     = errors.New("error getting keys from database")
)

// KGS is the core for Key Generation Service.
type KGS struct {
	db repository.KGSDatabase
}

// New creates a new instance of KGS and generate keys concurrently to the database.
func New(db repository.KGSDatabase, defaultPoolSize int, keyLength int) (*KGS, error) {
	// PostgreSQL has a default limit of 115 concurrent connections.
	// If connection(read/write goroutines) exceeded the limit,
	// it triggers the "FATAL: sorry, too many clients already" error, causing incoming connections to be rejected.
	maxDatabaseConnections := 100

	kgs := &KGS{db: db}
	ctx := context.TODO()

	errChan := make(chan error)
	doneChan := make(chan struct{})

	// Buffered semaphoreChan blocks goroutine from starting when channel is full.
	// Acts as a pool that allows token to be acquired(put token in semaphore) or to be released(drain semaphore).
	semaphoreChan := make(chan struct{}, maxDatabaseConnections)

	var wg sync.WaitGroup
	for i := 0; i < defaultPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Put token to semaphore when start goroutine.
			semaphoreChan <- struct{}{}
			defer func() {
				// Release semaphore(allowing other goroutines to put a new token to semaphore) after function done.
				<-semaphoreChan
			}()

			for {
				key, err := generateKey(keyLength)
				if err != nil {
					errChan <- ErrInvalidKeyLength
					return
				}
				exist, err := kgs.db.KeyExist(ctx, key)
				if err != nil && !errors.Is(err, repository.ErrKeyNotFound) {
					// send error to errChan, then shut the program
					log.Println(err)
					errChan <- ErrRepoError
					return
				}

				if !exist {
					err = kgs.db.WriteKey(ctx, key)
					if err != nil {
						// send error to errChan, then shut the program
						log.Println(err)
						errChan <- ErrRepoError
						return
					}
					break
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:
		return kgs, nil
	}
}

// generateKey generates the last four char in the shortenURL.
// Shortened URL should be in form: 'https://goShorten/1234'.
// Size of shortened URL is 16 bytes. Key in form '****' should have 62 ^ 4 variations.
// All variations take 16 * 62 ^ 4 bytes = 2,3642,1376 bytes, which is less than 0.3 GB.
// Implement a Key Generation Service using 1GB storage should be enough.
func generateKey(length int) (string, error) {
	if length <= 0 {
		return "", ErrInvalidKeyLength
	}
	res := make([]byte, length)

	for i := 0; i < length; i++ {
		res[i] = LetterBytes[rand.Intn(len(LetterBytes))]
	}
	return string(res), nil
}

// GetKeys fetches an array of keys with length requiredKeys from the Key Generation Service database.
func (k *KGS) GetKeys(ctx context.Context, requiredKeys int) ([]string, error) {
	ctrlCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	keys, err := k.db.GetKeys(ctrlCtx, requiredKeys)
	if err != nil {
		// TODO: Log the error
		return nil, ErrGetKeysError
	}

	return keys, nil
}
