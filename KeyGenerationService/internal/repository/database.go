package repository

// KGSDatabase is the interface that wraps writing and fetching keys from a Key Generation Service Database.
type KGSDatabase interface {
	KeyExist(string) (bool, error)
	WriteKey(string) error
	GetKeys(int) ([]string, error)
}
