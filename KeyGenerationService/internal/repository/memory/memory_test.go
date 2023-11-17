package memory

import (
	"KeyGenerationService/internal/repository"
	"context"
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Errorf("Error creating a new in-memory database: %v.\n", err)
	}
	if db == nil {
		t.Errorf("Error shouldn't have nil in-memory database.")
	}
}

func TestInMemoryDB_KeyExist(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Errorf("Error creating a new in-memory database: %v.\n", err)
	}
	ctx := context.Background()

	testKey := "1234"

	// 1. Fetch key that doesn't exist.
	ok, err := db.KeyExist(ctx, testKey)
	if !errors.Is(err, repository.ErrKeyNotFound) {
		t.Errorf("Error wrong error: Have %v, want %v.\n", err, repository.ErrKeyNotFound)
	}

	if ok {
		t.Errorf("Error shouldn't have fetched key that doesn't exist.\n")
	}

	// 2. Store key in keys, check key existence.
	db.Keys.Store(testKey, struct{}{})
	ok, err = db.KeyExist(ctx, testKey)
	if err != nil && !ok {
		t.Errorf("Error checking key existence: %v.\n", err)
	}

	db.Keys.Delete(testKey)

	// 3. Store key in UsedKeys, check key existence.
	db.UsedKeys.Store(testKey, struct{}{})
	ok, err = db.KeyExist(ctx, testKey)
	if err != nil && !ok {
		t.Errorf("Error checking key existence: %v.\n", err)
	}

	db.UsedKeys.Delete(testKey)
}

func TestInMemoryDB_WriteKey(t *testing.T) {
	inMemory, err := New()
	if err != nil {
		t.Errorf("Error creating a new in-memory database: %v.\n", err)
	}
	ctx := context.Background()

	testKey := "1234"

	err = inMemory.WriteKey(ctx, testKey)
	if err != nil {
		t.Errorf("Error writing key to in-memory database: %v.\n", err)
	}

	_, ok := inMemory.Keys.Load(testKey)
	if !ok {
		t.Errorf("Error written key is not in in-memory database.\n")
	}

}

func TestInMemoryDB_GetKeys(t *testing.T) {
	inMemory, err := New()
	if err != nil {
		t.Errorf("Error creating a new in-memory database: %v.\n", err)
	}
	ctx := context.Background()

	testKeys := []string{"0123", "1234", "2345", "3456", "4567", "5678", "6789", "7890"}
	for _, key := range testKeys {
		inMemory.Keys.Store(key, struct{}{})
	}

	// Test different requiredKeys.
	// TODO: Test when key required keys is greater than the amount of valid keys left.
	cases := []int{-1, 0, 1, 3, 10}

	for _, requiredKeys := range cases {
		result, err := inMemory.GetKeys(ctx, requiredKeys)
		if err != nil {
			if (requiredKeys <= 0 || requiredKeys > len(testKeys)) && !errors.Is(err, repository.ErrKeyOOR) {
				t.Errorf("Error incorrect error: Have %v, want %v.\n", err, repository.ErrKeyOOR)
			}
		} else {
			if len(result) != requiredKeys {
				t.Errorf("Error get keys error: Invalid keys length.\n")
			}
		}
	}
}
