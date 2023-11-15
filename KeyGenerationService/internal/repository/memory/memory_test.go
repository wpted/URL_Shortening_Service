package memory

import (
	"KeyGenerationService/internal/repository"
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Errorf("Error shouldn't have error on creating a new in-memory database")
	}
	if db == nil {
		t.Errorf("Error shouldn't have nil in-memory database.")
	}
}

func TestInMemoryDB_KeyExist(t *testing.T) {
	inMemory, err := New()
	if err != nil {
		t.Errorf("Error initializing in-memory database: %v.\n", err)
	}

	testKey := "1234"

	// 1. Fetch key that doesn't exist.
	ok, err := inMemory.KeyExist(testKey)
	if err != nil {
		t.Errorf("Error checking existence: %v\n.", err)
	}

	if ok {
		t.Errorf("Error fetched unwanted key.\n")
	}

	// 2. Store key in keys, check key existence.
	inMemory.Keys.Store(testKey, struct{}{})
	ok, err = inMemory.KeyExist(testKey)
	if err != nil {
		t.Errorf("Error checking existence: %v\n.", err)
	}

	if !ok {
		t.Errorf("Error key should exist.\n")
	}
	inMemory.Keys.Delete(testKey)

	// 3. Store key in UsedKeys, check key existence.
	inMemory.UsedKeys.Store(testKey, struct{}{})
	ok, err = inMemory.KeyExist(testKey)
	if err != nil {
		t.Errorf("Error checking existence: %v\n.", err)
	}

	if !ok {
		t.Errorf("Error key should exist.\n")
	}
	inMemory.UsedKeys.Delete(testKey)
}

func TestInMemoryDB_WriteKey(t *testing.T) {
	inMemory, err := New()
	if err != nil {
		t.Errorf("Error initializing in-memory database: %v.\n", err)
	}

	testKey := "1234"

	err = inMemory.WriteKey(testKey)
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
		t.Errorf("Error initializing in-memory database: %v.\n", err)
	}

	testKeys := []string{"0123", "1234", "2345", "3456", "4567", "5678", "6789", "7890"}
	for _, key := range testKeys {
		inMemory.Keys.Store(key, struct{}{})
	}

	// Test different requiredKeys.
	// TODO: Test when key required keys is greater than the amount of valid keys left.
	cases := []int{-1, 0, 1, 3, 10}

	for _, requiredKeys := range cases {
		result, err := inMemory.GetKeys(requiredKeys)
		if err != nil {
			if !errors.Is(err, repository.ErrNegativeKey) {
				t.Errorf("Error incorrect error: Have %v, want %v.\n", err, repository.ErrNegativeKey)
			}
		} else {
			if len(result) != requiredKeys {
				t.Errorf("Error get keys error: Invalid keys length.\n")
			}
		}
	}
}
