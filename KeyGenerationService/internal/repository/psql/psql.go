package psql

import (
	"KeyGenerationService/internal/repository"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

// DB used for Key Generation Service.
type DB struct {
	db *sql.DB
}

// New creates a new instance of DB.
func New(user, password, database string) (*DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, database))
	if err != nil {
		return nil, repository.ErrDatabaseError
	}

	return &DB{db: db}, nil
}

// KeyExist checks whether a key exist within DB.
func (d *DB) KeyExist(key string) (bool, error) {
	var value string
	inKeys, inUsedKeys := true, true

	// Check key existence in keys.
	query := "SELECT values FROM keys WHERE values=$1"
	row := d.db.QueryRow(query, key)

	err := row.Scan(&value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, repository.ErrDatabaseError
		}
		inKeys = false
	}

	// Check key existence in used_keys.
	query = "SELECT values FROM used_keys WHERE values=$1"
	row = d.db.QueryRow(query, key)

	err = row.Scan(&value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return false, repository.ErrDatabaseError
		}
		inUsedKeys = false
	}

	if inKeys || inUsedKeys {
		return true, nil
	} else {
		return false, repository.ErrKeyNotFound
	}
}

// WriteKey stores the given key to DB.
func (d *DB) WriteKey(key string) error {
	query := "INSERT INTO keys(values) VALUES($1)"
	_, err := d.db.Exec(query, key)
	if err != nil {
		return repository.ErrDatabaseError
	}

	return nil
}

// GetKeys fetches an array of keys.
// The fetched keys are considered used and will be moved to used_keys for further usage.
func (d *DB) GetKeys(requiredKeys int) ([]string, error) {
	// Cannot have negative or zero requiredKeys.
	if requiredKeys <= 0 {
		return []string{}, repository.ErrNegativeKey
	}

	// Cannot have requiredKeys greater than what we have in 'keys'.
	// This shouldn't happen since we always assume that we have enough keys in pool waiting.
	rows, err := d.db.Query("SELECT COUNT(*) FROM keys;")
	if err != nil {
		return nil, repository.ErrDatabaseError
	}

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return nil, repository.ErrDatabaseError
		}
	}
	_ = rows.Close()
	if requiredKeys > count {
		return []string{}, repository.ErrKeyOutOfRange
	}

	// Create an array that stores all fetched keys.
	result := make([]string, requiredKeys)
	query := "SELECT values FROM keys FETCH FIRST $1 ROWS ONLY"
	rows, err = d.db.Query(query, requiredKeys)
	if err != nil {
		return nil, repository.ErrDatabaseError
	}

	i := 0
	for rows.Next() {
		var key string
		err := rows.Scan(&key)
		if err != nil {
			return nil, repository.ErrDatabaseError
		}
		result[i] = key

		_, err = d.db.Exec("DELETE FROM keys WHERE values=$1", key)
		if err != nil {
			return nil, repository.ErrDatabaseError
		}
		_, err = d.db.Exec("INSERT INTO used_keys(values) VALUES($1)", key)
		if err != nil {
			return nil, repository.ErrDatabaseError
		}
		i++
	}

	_ = rows.Close()

	return result, nil
}
