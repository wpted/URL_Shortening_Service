package controller

import (
	"KeyGenerationService/internal/repository/memory"
	"KeyGenerationService/internal/repository/psql"
	"context"
	"errors"
	"testing"
)

func TestNew(t *testing.T) {

	t.Run("Test in-memory database", func(t *testing.T) {
		db, err := memory.New()
		if err != nil {
			t.Errorf("Error creating instance DB.\n")
		}
		if err != nil {
			t.Errorf("Error creating database: %v.\n", err)
		}

		kgs, err := New(db, 100000, 4)
		if err != nil || kgs == nil {
			t.Errorf("Error creating controller: %v.\n", err)
		}
	})

	t.Run("Test relational database", func(t *testing.T) {
		db, err := psql.New("URLShortenerUser", "URLShortenerPassword", "KeyGenerationService")
		if err != nil {
			t.Errorf("Error creating instance DB.\n")
		}
		if err != nil {
			t.Errorf("Error creating database: %v.\n", err)
		}

		defaultPoolSizeCases := []int{-1, 0, 1000}

		for _, defaultPoolSize := range defaultPoolSizeCases {
			kgs, err := New(db, defaultPoolSize, 4)
			if err != nil {
				if defaultPoolSize < 0 && !errors.Is(err, ErrInvalidPoolSize) {
					t.Errorf("Error incorrect error: Have %v, want %v.\n", err, ErrInvalidPoolSize)
				}
				if defaultPoolSize >= 0 && kgs == nil {
					t.Errorf("Error creating controller: %v.\n", err)
				}
			}
		}

		db.CleanUp()
	})

}

func TestKGS_GetKeys(t *testing.T) {
	ctx := context.Background()

	db, err := memory.New()
	if err != nil {
		t.Errorf("Error creating instance DB.\n")
	}

	defaultPoolSize := 100
	kgs, err := New(db, defaultPoolSize, 4)
	if err != nil || kgs == nil {
		t.Errorf("Error creating controller: %v.\n", err)
	}

	requiredKeysCases := []int{-1, 0, 10, 101}
	for _, requiredKeys := range requiredKeysCases {
		keys, err := kgs.GetKeys(ctx, requiredKeys)
		if err != nil {
			var ctrlError *KGSError
			if (requiredKeys <= 0 || requiredKeys > defaultPoolSize) && !errors.As(err, &ctrlError) {
				t.Errorf("Error incorrect error: %v.\n", err)
			} else if 0 < requiredKeys && requiredKeys <= defaultPoolSize {
				t.Errorf("Error getting keys from database: %v.\n", err)
			}
		} else {
			if len(keys) != requiredKeys {
				t.Errorf("Error fetched keys length is incorrect: Have %v, want %v.\n", len(keys), requiredKeys)
			}
		}
	}
}

func Test_generateKey(t *testing.T) {
	lengthCases := []int{-1, 0, 4}

	for _, length := range lengthCases {
		key, err := generateKey(length)

		if err != nil {
			if len(key) <= 0 && !errors.Is(err, ErrInvalidKeyLength) {
				t.Errorf("Error incorrect error: Have %v, want %v.\n", err, ErrInvalidKeyLength)
			}
			if len(key) > 0 {
				t.Errorf("Error generating key: %v.\n", err)
			}
		} else {
			if len(key) != length {
				t.Errorf("Error incorrect generated key length: Have %v, want %v.\n", len(key), length)
			}
		}
	}
}
