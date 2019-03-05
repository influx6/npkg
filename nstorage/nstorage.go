package nstorage

import "time"

// ByteStore defines a storage interface defining what we expect to
// be provided for storing a byte slice with a underline key.
type ByteStore interface {
	Keys() ([]string, error)
	Save(string, []byte) error
	Get(string) ([]byte, error)
	Exists(string) (bool, error)
	Update(string, []byte) error
	Remove(string) ([]byte, error)
	Each(fn func([]byte, string) bool) error
}

// ExpirableStore composes the ByteStore providing the
// expiration versions of Save and TTL extension.
type ExpirableStore interface {
	ByteStore

	// TTL should return current expiry value of
	// giving key in millisecond.
	TTL(string) (time.Duration, error)

	// ExtendTTL should extend expiry by giving duration.
	// Add new duration to the remaining ttl of key.
	ExtendTTL(string, time.Duration) error

	// SaveTTL save giving key with giving expiry.
	SaveTTL(string, []byte, time.Duration) error

	// UpdateTTL updates giving key with giving expiry.
	// It should update key's value and add giving duration
	// to remaining time of key.
	UpdateTTL(string, []byte, time.Duration) error
}
