package nbadger

import (
	"time"
	"unsafe"

	"github.com/dgraph-io/badger"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nstorage"
)

var _ nstorage.ExpirableStore = (*BadgerStore)(nil)
var _ nstorage.QueryableByteStore = (*BadgerStore)(nil)

// BadgerStore implements session management, storage and access using Badger as
// underline store.
type BadgerStore struct {
	prefix string
	ops    badger.Options
	iter   badger.IteratorOptions
	Db     *badger.DB
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
	rd.Db = db
	return nil
}

// Keys returns all giving keys of elements within store.
func (rd *BadgerStore) Keys() ([]string, error) {
	var keys = make([]string, 0)
	var err = rd.Db.View(func(txn *badger.Txn) error {
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

// Find returns all matching results within using giving functions.
func (rd *BadgerStore) Find(fn func([]byte, string) bool) error {
	var err = rd.Db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		if rd.prefix == "" {
			for iterator.Rewind(); iterator.Valid(); iterator.Next() {
				var item = iterator.Item()
				if item.IsDeletedOrExpired() {
					continue
				}
				var stop = false
				var err = item.Value(func(value []byte) error {
					if !fn(value, string(item.Key())) {
						stop = true
					}
					return nil
				})

				if err != nil {
					return nerror.WrapOnly(err)
				}

				if stop {
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
			var stop = false
			var err = item.Value(func(value []byte) error {
				if !fn(value, string(item.Key())) {
					stop = true
				}
				return nil
			})

			if err != nil {
				return nerror.WrapOnly(err)
			}

			if stop {
				return nil
			}
		}
		return nil
	})
	return err
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
	return rd.Db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		if rd.prefix == "" {
			for iterator.Rewind(); iterator.Valid(); iterator.Next() {
				var item = iterator.Item()
				if item.IsDeletedOrExpired() {
					continue
				}
				var stop = false
				var err = item.Value(func(value []byte) error {
					if !fn(value, string(item.Key())) {
						stop = true
					}
					return nil
				})

				if err != nil {
					return nerror.WrapOnly(err)
				}

				if stop {
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
			var stop = false
			var err = item.Value(func(value []byte) error {
				if !fn(value, string(item.Key())) {
					stop = true
				}
				return nil
			})

			if err != nil {
				return nerror.WrapOnly(err)
			}

			if stop {
				return nil
			}
		}
		return nil
	})
}

// Exists returns true/false if giving key exists.
func (rd *BadgerStore) Exists(key string) (bool, error) {
	var exist bool
	if err := rd.Db.View(func(txn *badger.Txn) error {
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
	if err := rd.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.WrapOnly(err)
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found")
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return nerror.WrapOnly(err)
		}

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
	return rd.Db.Update(func(txn *badger.Txn) error {
		var op badger.Entry
		op.Value = data
		op.Key = string2Bytes(key)
		op.ExpiresAt = uint64(time.Now().Add(expiration).Unix())

		if err := txn.SetEntry(&op); err != nil {
			return nerror.WrapOnly(err)
		}
		return nil
	})
}

// TTL returns giving expiration time for giving key.
func (rd *BadgerStore) TTL(key string) (time.Duration, error) {
	var ttl time.Duration
	var err = rd.Db.Update(func(txn *badger.Txn) error {
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
	return rd.Db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var op badger.Entry
		op.Value = value
		op.Key = string2Bytes(key)

		var expr = ttlDur(item.ExpiresAt(), 0)
		op.ExpiresAt = uint64(time.Now().Add(expr + expiration).Unix())

		if err := txn.SetEntry(&op); err != nil {
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
	return rd.Db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var op badger.Entry
		op.Value = value
		op.Key = string2Bytes(key)
		op.ExpiresAt = uint64(time.Now().Add(expiration).Unix())

		if err := txn.SetEntry(&op); err != nil {
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
	return rd.Db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var op badger.Entry
		op.Value = value
		op.Key = string2Bytes(key)

		var ttl time.Time
		if expiration == 0 {
			ttl = time.Unix(int64(item.ExpiresAt()), 0)
		} else {
			ttl = time.Now().Add(expiration)
		}

		op.ExpiresAt = uint64(ttl.Unix())
		if err := txn.SetEntry(&op); err != nil {
			return err
		}
		return nil
	})
}

// Remove removes underline key from the Badger store after retrieving it and
// returning giving session.
func (rd *BadgerStore) Remove(key string) ([]byte, error) {
	var old []byte
	err := rd.Db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(string2Bytes(key))
		if err != nil {
			return err
		}

		old = make([]byte, item.ValueSize())
		old, err = item.ValueCopy(old)
		if err != nil {
			return err
		}

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
