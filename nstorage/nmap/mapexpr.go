package nmap

import (
	"time"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nstorage"
)

var _ nstorage.ExpirableStore = (*ExprByteStore)(nil)
var _ nstorage.QueryableByteStore = (*ExprByteStore)(nil)

// ExprByteStore implements an expiring byte store that
// matches the nstorage.ExpirableStorage interface.
type ExprByteStore struct {
	cache *ExpiringByteMap
}

// NewExprByteStore returns a new instance of a ExprByteStore.
func NewExprByteStore(initial ...uint) *ExprByteStore {
	var expr ExprByteStore
	expr.cache = NewExpiringByteMap(initial...)
	return &expr
}

// Keys returns the list of all keys registered to giving store.
func (expr *ExprByteStore) Keys() ([]string, error) {
	var keys []string
	expr.cache.GetMany(func(values map[string]ExpiringValue) {
		keys = make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
	})
	return keys, nil
}

// Get returns the giving value of key if it exists and has not expired.
func (expr *ExprByteStore) Get(k string) ([]byte, error) {
	if !expr.cache.Has(k) {
		return nil, nerror.New("not found")
	}
	return expr.cache.Get(k), nil
}

// TTL returns the remaining time left for giving key before expiration.
//
// A zero value means has no expiration.
func (expr *ExprByteStore) TTL(k string) (time.Duration, error) {
	return expr.cache.TTL(k), nil
}

// Exists returns true/false if giving key exists.
func (expr *ExprByteStore) Exists(k string) (bool, error) {
	return expr.cache.Has(k), nil
}

// Save adds giving key and value into store.
func (expr *ExprByteStore) Save(k string, v []byte) error {
	expr.cache.Set(k, v, 0)
	return nil
}

// ExtendTTL extends expiration or sets expiration of giving key if it has not
// expired yet.
//
// A expiration value of zero means to persist the giving key.
func (expr *ExprByteStore) ExtendTTL(k string, t time.Duration) error {
	expr.cache.ExtendTTL(k, t)
	return nil
}

// ResetTTL resets expiration or sets expiration of giving key if it has not
// expired yet.
//
// A expiration value of zero means to persist the giving key.
func (expr *ExprByteStore) ResetTTL(k string, t time.Duration) error {
	expr.cache.ResetTTL(k, t)
	return nil
}

// Updates updates giving key and value into store.
func (expr *ExprByteStore) Update(k string, v []byte) error {
	expr.cache.Set(k, v, 0)
	return nil
}

// SaveTTL updates giving key and value into store with expiration value.
func (expr *ExprByteStore) SaveTTL(k string, v []byte, t time.Duration) error {
	expr.cache.Set(k, v, t)
	return nil
}

// UpdateTTL updates giving key and value into store with expiration value.
func (expr *ExprByteStore) UpdateTTL(k string, v []byte, t time.Duration) error {
	if !expr.cache.Has(k) {
		return nerror.New("key does not exists")
	}

	expr.cache.Set(k, v, t)
	return nil
}

// Each alternatives through all keys and values from underline cache.
//
// To ensure no-undesired behaviour, ensure to copy the value to avoid
// possible change to it, as the underline store owns the giving value
// slice and maybe re-used as it sees fit.
func (expr *ExprByteStore) Each(fn func([]byte, string) bool) error {
	expr.cache.GetMany(func(values map[string]ExpiringValue) {
		for key, value := range values {
			if !fn(value.Value, key) {
				return
			}
		}
	})
	return nil
}

// Find returns all elements matching giving function and count.
func (expr *ExprByteStore) Find(fn func([]byte, string) bool) error {
	expr.cache.GetMany(func(values map[string]ExpiringValue) {
		for key, value := range values {
			if !fn(value.Value, key) {
				return
			}
		}
	})
	return nil
}

// Remove deletes giving key from underling store.
func (expr *ExprByteStore) Remove(k string) ([]byte, error) {
	var v []byte
	var found bool
	expr.cache.GetMany(func(values map[string]ExpiringValue) {
		for key, value := range values {
			if key == key {
				found = true
				v = value.Value
				delete(values, key)
				return
			}
		}
	})

	if !found {
		return nil, nerror.New("Key does not exists")
	}
	return v, nil
}
