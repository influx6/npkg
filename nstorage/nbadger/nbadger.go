package nbadger

import (
	"time"
	"unsafe"

	"github.com/dgraph-io/badger"
	"github.com/gokit/npkg/nerror"
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
		return err
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
				var mycopy = make([]byte, len(item.Key()))
				copy(mycopy, item.Key())
				keys = append(keys, bytes2String(mycopy))
			}
			return nil
		}

		var prefix = []byte(rd.prefix)
		for iterator.Rewind(); iterator.ValidForPrefix(prefix); iterator.Next() {
			var item = iterator.Item()
			var mycopy = make([]byte, len(item.Key()))
			copy(mycopy, item.Key())
			keys = append(keys, bytes2String(mycopy))
		}
		return nil
	})
	return keys, err
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
				var value, err = item.Value()
				if err != nil {
					return err
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
			var value, err = item.Value()
			if err != nil {
				return err
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
			return err
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
			return err
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found")
		}

		dbValue, err := item.Value()
		if err != nil {
			return err
		}

		value = make([]byte, len(dbValue))
		copy(value, dbValue)
		return nil
	}); err != nil {
		return nil, err
	}
	return value, nil
}

// Save adds giving session into storage using Badger as underline store.
func (rd *BadgerStore) Save(key string, data []byte, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		if expiration > 0 {
			if err := txn.SetWithTTL(string2Bytes(key), data, expiration); err != nil {
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

// Update updates giving session stored with giving key. It updates
// the underline data.
func (rd *BadgerStore) Update(key string, data []byte, expiration time.Duration) error {
	return rd.db.Update(func(txn *badger.Txn) error {
		if expiration > 0 {
			if err := txn.SetWithTTL(string2Bytes(key), data, expiration); err != nil {
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

		old = make([]byte, len(value))
		copy(old, value)
		return txn.Delete(string2Bytes(key))
	})
	return old, err
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
