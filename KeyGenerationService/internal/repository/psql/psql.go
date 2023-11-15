package psql

import (
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
		// TODO: Handle database connection error
		return nil, err
	}

	return &DB{db: db}, nil
}

// KeyExist checks whether a key exist within DB.
func (d *DB) KeyExist(key string) (bool, error) {
	var value string
	query := "SELECT values FROM keys WHERE values=$1"
	row := d.db.QueryRow(query, key)

	err := row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Key doesn't exist.
			return false, nil
		} else {
			// TODO: Handle read error.
			return false, err
		}
	}
	return true, nil
}

// WriteKey stores the given key to DB.
func (d *DB) WriteKey(key string) error {
	query := "INSERT INTO keys(values) VALUES($1)"
	_, err := d.db.Exec(query, key)
	if err != nil {
		// TODO: Handle insert error.
	}

	return nil
}

// GetKeys fetches an array of keys.
// The fetched keys are considered used and will be moved to used_keys for further usage.
func (d *DB) GetKeys(requiredKeys int) ([]string, error) {
	result := make([]string, requiredKeys)
	query := "SELECT values FROM keys LIMIT $1"
	rows, err := d.db.Query(query, requiredKeys)
	if err != nil {
		// TODO: Handle read error.
		return nil, err
	}

	i := 0
	for rows.Next() {
		var key string
		err := rows.Scan(&key)
		if err != nil {
			// TODO: Handle scan error.
		}
		result[i] = key

		_, err = d.db.Exec("DELETE FROM keys WHERE values=$1", key)
		// TODO: Handle delete error.
		_, err = d.db.Exec("INSERT INTO used_keys(values) VALUES($1)", key)
		// TODO: Handle insert error.
		i++
	}

	return result, nil
}
