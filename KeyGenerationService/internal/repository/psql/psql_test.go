package psql

import (
	"KeyGenerationService/internal/repository"
	"context"
	"database/sql"
	"errors"
	"testing"
)

// Before testing, be sure the database for testing is connected.
// TODO: Create database for testing.

func TestNew(t *testing.T) {
	// TODO: Replace the credentials.
	db, err := New("URLShortenerUser", "URLShortenerPassword", "KeyGenerationService")
	if err != nil {
		t.Errorf("Error shouldn't have an error when creating instance DB.\n")
	}

	if db == nil {
		t.Errorf("Error shouldn't return a nil DB instance.\n")
	}
}

func TestDB_KeyExist(t *testing.T) {
	db, err := New("URLShortenerUser", "URLShortenerPassword", "KeyGenerationService")
	if err != nil {
		t.Errorf("Error shouldn't have an error when creating instance DB.\n")
	}

	testKey := "test_key"

	ctx := context.Background()
	// 1. Fetch key that doesn't exist.
	ok, err := db.KeyExist(ctx, testKey)
	if err != nil && !errors.Is(err, repository.ErrKeyNotFound) {
		t.Errorf("Error wrong error: Have %v, want %v.\n", err, repository.ErrKeyNotFound)
	}

	if ok {
		t.Errorf("Error shouldn't have fetched key that doesn't exist.\n")
	}

	// 2. Store key in keys, check key existence.
	_, _ = db.db.Exec("INSERT INTO keys(values) VALUES ($1)", testKey)
	ok, err = db.KeyExist(ctx, testKey)
	if err != nil && !errors.Is(err, repository.ErrKeyNotFound) {
		// Database error.
		t.Errorf("Error checking key existence: %v.\n", err)
	}
	if !ok {
		// Key not found.
		t.Errorf("Error checking key existence: %v.\n", err)
	}
	_, _ = db.db.Exec("DELETE FROM keys WHERE values = $1", testKey)

	// 3. Store key in UsedKeys, check key existence.
	_, _ = db.db.Exec("INSERT INTO used_keys(values) VALUES ($1)", testKey)
	ok, err = db.KeyExist(ctx, testKey)
	if err != nil && !errors.Is(err, repository.ErrKeyNotFound) {
		// Database error.
		t.Errorf("Error checking key existence: %v.\n", err)
	}
	if !ok {
		// Key not found.
		t.Errorf("Error checking key existence: %v.\n", err)
	}
	_, _ = db.db.Exec("DELETE FROM used_keys WHERE values = $1", testKey)
}

func TestDB_WriteKey(t *testing.T) {
	db, err := New("URLShortenerUser", "URLShortenerPassword", "KeyGenerationService")
	if err != nil {
		t.Errorf("Error shouldn't have an error when creating instance DB.\n")
	}
	ctx := context.Background()

	testKey := "test_key"
	err = db.WriteKey(ctx, testKey)
	if err != nil {
		t.Errorf("Error writing key to database: %v.\n", err)
	}

	var haveKey string
	query := "SELECT values FROM keys WHERE values = $1"
	row := db.db.QueryRow(query, testKey)
	err = row.Scan(&haveKey)
	if errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Error written key is not in database.\n")
	}
	if haveKey != testKey {
		t.Errorf("Error fetched wrong key: Have %v, want %v\n.", haveKey, testKey)
	}

	_, _ = db.db.Exec("DELETE FROM keys where values = $1", testKey)
}

func TestDB_GetKeys(t *testing.T) {
	db, err := New("URLShortenerUser", "URLShortenerPassword", "KeyGenerationService")
	if err != nil {
		t.Errorf("Error shouldn't have an error when creating instance DB.\n")
	}
	testKeys := []string{"test_key1", "test_key2", "test_key3", "test_key4", "test_key5", "test_key6"}

	for _, testKey := range testKeys {
		_, _ = db.db.Exec("INSERT INTO keys(values) VALUES ($1)", testKey)
	}

	ctx := context.Background()

	// Shouldn't have requiredKeysCases greater than existing testKeys. (Should check if there's enough key and update periodically.)
	requiredKeysCases := []int{-1, 0, 1, 3, 10}
	for _, requiredKeys := range requiredKeysCases {
		keys, err := db.GetKeys(ctx, requiredKeys)
		if err != nil {
			if requiredKeys <= 0 && !errors.Is(err, repository.ErrNegativeKey) {
				t.Errorf("Error incorrect error: Have %v, want %v.\n", err, repository.ErrNegativeKey)
			}
			if requiredKeys > len(testKeys) && !errors.Is(err, repository.ErrKeyOutOfRange) {
				t.Errorf("Error incorrect error: Have %v, want %v.\n", err, repository.ErrKeyOutOfRange)
			}
		} else {
			if len(keys) != requiredKeys {
				t.Errorf("Error get keys error: Invalid keys length.\n")
			}
		}
	}

	// Clean the table.
	_, _ = db.db.Exec("DELETE FROM keys")
}
