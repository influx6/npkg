package nbadger

import (
	"time"
	"unsafe"

	"github.com/dgraph-io/badger"

	"github.com/influx6/npkg/nerror"
)

// BadgerStore implements session management, storage and access using Badger as
// underline store.
type BadgerStore struct {
	prefix string
	ops    badger.Options
	iter   badger.IteratorOptions
	db     *badger.DB
}

// NewBadgerStore returns a new instance of a Badger store using provided prefix if present.
func NewBadgerStore(prefix string, ops badger.Options, iterator badger.IteratorOptions) (*BadgerStore, error) {
	var red BadgerStore
	red.ops = ops
	red.prefix = prefix
	red.iter = iterator
	if err := red.createConnection(); err != nil {
		return nil, err
	}
	return &red, nil
}

// createConnection attempts to create a new Badger connection.
func (rd *BadgerStore) createConnection() error {
	db, err := badger.Open(rd.ops)
	if err != nil {
		return nerror.WrapOnly(err)
	}
	rd.db = db
	return nil
}

// Keys returns all giving keys of elements within store.
func (rd *BadgerStore) Keys() ([]string, error) {
	var keys = make([]string, 0)
	var err = rd.db.View(func(txn *badger.Txn) error {
		var newIterator = rd.iter
		newIterator.PrefetchValues = false
		var iterator = txn.NewIterator(newIterator)
		defer iterator.Close()

		if rd.prefix == "" {
			for iterator.Rewind(); iterator.Valid(); iterator.Next() {
				var item = iterator.Item()
				if item.IsDeletedOrExpired() {
					continue
				}
				keys = append(keys, bytes2String(copyBytes(item.Key())))
			}
			return nil
		}

		var prefix = []byte(rd.prefix)
		for iterator.Rewind(); iterator.ValidForPrefix(prefix); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			keys = append(keys, bytes2String(copyBytes(item.Key())))
		}
		return nil
	})
	return keys, err
}

// FindAll returns all match elements for giving function.
func (rd *BadgerStore) FindAll(fn func([]byte, string) bool, count int) ([][]byte, error) {
	return rd.FindEach(fn, -1)
}

// Find returns the single result matching giving function.
func (rd *BadgerStore) Find(fn func([]byte, string) bool) ([]byte, error) {
	var res, err = rd.FindEach(fn, 1)
	if err != nil {
		return nil, err
	}
	if len(res) == 1 {
		return res[0], err
	}
	return nil, nil
}

// FindEach returns all matching results within count if not -1 using giving functions.
func (rd *BadgerStore) FindEach(fn func([]byte, string) bool, count int) ([][]byte, error) {
	var results [][]byte
	var err = rd.db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		if rd.prefix == "" {
			for iterator.Rewind(); iterator.Valid(); iterator.Next() {
				if count > 0 && len(results) == count {
					return nil
				}
				var item = iterator.Item()
				if item.IsDeletedOrExpired() {
					continue
				}
				var value, err = item.Value()
				if err != nil {
					return nerror.WrapOnly(err)
				}

				if fn(value, string(item.Key())) {
					results = append(results, value)
					continue
				}
			}
			return nil
		}

		var prefix = []byte(rd.prefix)
		for iterator.Rewind(); iterator.ValidForPrefix(prefix); iterator.Next() {
			if count > 0 && len(results) == count {
				return nil
			}
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			var value, err = item.Value()
			if err != nil {
				return nerror.WrapOnly(err)
			}
			if fn(value, string(item.Key())) {
				results = append(results, value)
				continue
			}
		}
		return nil
	})
	return results, err
}

// Each runs through all elements for giving store, skipping keys
// in Badger who have no data or an empty byte slice.
//
// Each byte slice provided is only valid for the call of
// the function, after which it becomes invalid as it can
// be re-used for efficient memory management, so ensure to copy
// given byte slice yourself within function to protect against
// undefined behaviour.
func (rd *BadgerStore) Each(fn func([]byte, string) bool) error {
	return rd.db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		if rd.prefix == "" {
			for iterator.Rewind(); iterator.Valid(); iterator.Next() {
				var item = iterator.Item()
				if item.IsDeletedOrExpired() {
					continue
				}
				var value, err = item.Value()
				if err != nil {
					return nerror.WrapOnly(err)
				}

				if !fn(value, string(item.Key())) {
					return nil
				}
			}
			return nil
		}

		var prefix = []byte(rd.prefix)
		for iterator.Rewind(); iterator.ValidForPrefix(prefix); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			var value, err = item.Value()
			if err != nil {
				return nerror.WrapOnly(err)
			}
			if !fn(value, string(item.Key())) {
				return nil
			}
		}
		return nil
	})
}

// Exists returns true/false if giving key exists.
func (rd *BadgerStore) Exists(key string) (bool, error) {
	var exist bool
	if err := rd.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.WrapOnly(err)
		}
		if item.IsDeletedOrExpired() {
			return nil
		}
		exist = true
		return nil
	}); err != nil {
		return false, err
	}
	return exist, nil
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *BadgerStore) Get(key string) ([]byte, error) {
	var value []byte
	if err := rd.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.WrapOnly(err)
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found")
		}

		dbValue, err := item.Value()
		if err != nil {
			return nerror.WrapOnly(err)
		}

		value = copyBytes(dbValue)
		return nil
	}); err != nil {
		return nil, err
	}
	return value, nil
}

// Save adds giving session into storage using Badger as underline store.
func (rd *BadgerStore) Save(key string, data []byte) error {
	return rd.SaveTTL(key, data, 0)
}

// SaveTTL adds giving session into storage using Badger as underline store, with provided
// expiration.
// Duration of 0 means no expiration.
func (rd *BadgerStore) SaveTTL(key string, data []byte, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		if expiration > 0 {
			if err := txn.SetWithTTL(string2Bytes(key), data, expiration); err != nil {
				return nerror.WrapOnly(err)
			}
			return nil
		}
		if err := txn.Set(string2Bytes(key), data); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})
}

// TTL returns giving expiration time for giving key.
func (rd *BadgerStore) TTL(key string) (time.Duration, error) {
	var ttl time.Duration
	var err = rd.db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		ttl = ttlDur(item.ExpiresAt(), 0)
		return nil
	})
	return ttl, err
}

// ExtendTTL resets new TTL for giving key if it has not expired and is still accessible.
//
// A expiration value of zero means to persist the giving key.
func (rd *BadgerStore) ExtendTTL(key string, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		value, err := item.Value()
		if err != nil {
			return err
		}

		if expiration == 0 {
			return txn.Set(string2Bytes(key), copyBytes(value))
		}

		var expr = ttlDur(item.ExpiresAt(), 0)
		var newExpr = expr + expiration

		if err := txn.SetWithTTL(string2Bytes(key), copyBytes(value), newExpr); err != nil {
			return err
		}
		return nil
	})
}

// ResetTTL resets expiration or sets expiration of giving key if it has not
// expired yet.
//
// A expiration value of zero means to persist the giving key.
func (rd *BadgerStore) ResetTTL(key string, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		value, err := item.Value()
		if err != nil {
			return err
		}

		if expiration == 0 {
			return txn.Set(string2Bytes(key), copyBytes(value))
		}

		if err := txn.SetWithTTL(string2Bytes(key), copyBytes(value), expiration); err != nil {
			return err
		}
		return nil
	})
}

// Update updates giving session stored with giving key. It updates
// the underline data for key, setting a 0 ttl duration.
func (rd *BadgerStore) Update(key string, data []byte) error {
	return rd.UpdateTTL(key, data, 0)
}

// UpdateTTL updates giving session stored with giving key. It updates
// the underline data.
// If a key has expired, and was deleted, then a error is returned.
//
// if expiration is zero then giving value expiration will not be reset but left
// as is.
func (rd *BadgerStore) UpdateTTL(key string, data []byte, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		if expiration > 0 {
			var item, err = txn.Get(string2Bytes(key))
			if err != nil {
				return nerror.Wrap(err, "Failed to retrieve key")
			}
			if item.IsDeletedOrExpired() {
				return nerror.New("not found, possibly expired")
			}

			var ttl = ttlDur(item.ExpiresAt(), 0) + expiration
			if err := txn.SetWithTTL(string2Bytes(key), data, ttl); err != nil {
				return err
			}
			return nil
		}
		if err := txn.Set(string2Bytes(key), data); err != nil {
			return err
		}
		return nil
	})
}

// Remove removes underline key from the Badger store after retrieving it and
// returning giving session.
func (rd *BadgerStore) Remove(key string) ([]byte, error) {
	var old []byte
	err := rd.db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return err
		}

		value, err := item.Value()
		if err != nil {
			return err
		}

		old = copyBytes(value)
		return txn.Delete(string2Bytes(key))
	})
	return old, err
}

//*****************************************************
// internal methods
//*****************************************************

func copyBytes(bu []byte) []byte {
	var cu = make([]byte, len(bu))
	copy(cu, bu)
	return cu
}

func ttlDur(tx uint64, ns int64) time.Duration {
	if tx == 0 {
		return 0
	}
	var ttl = time.Unix(int64(tx), ns)
	return time.Now().Sub(ttl)
}

//*****************************************************
// unsafe methods
//*****************************************************

// byte2String converts a byte slice into a string.
func bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

// string2Bytes converts a string into a byte slice.
func string2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}
