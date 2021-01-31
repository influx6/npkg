package nbadger

import (
	regexp2 "regexp"
	"time"

	"github.com/dgraph-io/badger/v2"

	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/nstorage"
	"github.com/influx6/npkg/nunsafe"
)

var _ nstorage.ExpirableStore = (*BadgerStore)(nil)

// BadgerStore implements session management, storage and access using Badger as
// underline store.
type BadgerStore struct {
	ops  badger.Options
	iter badger.IteratorOptions
	Db   *badger.DB
}

// NewBadgerStore returns a new instance of a Badger store using provided prefix if present.
func NewBadgerStore(ops badger.Options, iterator badger.IteratorOptions) (*BadgerStore, error) {
	var red BadgerStore
	red.ops = ops
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
		var iteratorOption = rd.iter
		iteratorOption.PrefetchValues = false
		var iterator = txn.NewIterator(iteratorOption)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			keys = append(keys, nunsafe.Bytes2String(copyBytes(item.Key())))
		}
		return nil
	})
	return keys, err
}

// ScanMatch uses the value of lastKey instead the index to allow scanning keys over
// a giving range, this is important and should be maintained and provided for
// this to work with BadgerDB.
func (rd *BadgerStore) ScanMatch(count int64, _ int64, lastKey string, regexp string) (nstorage.ScanResult, error) {
	if len(regexp) == 0 {
		regexp = ".+"
	}

	var regx, rgErr = regexp2.Compile(regexp)
	if rgErr != nil {
		return nstorage.ScanResult{}, nerror.WrapOnly(rgErr)
	}

	var isFinished bool
	var keys = make([]string, 0, 10)
	var err = rd.Db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		if len(lastKey) != 0 {
			iterator.Seek(nunsafe.String2Bytes(lastKey))
		} else {
			iterator.Rewind()
		}

		for ; iterator.Valid(); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			if regx != nil && !regx.Match(item.Key()) {
				continue
			}

			keys = append(keys, nunsafe.Bytes2String(item.Key()))
			if currCount := len(keys); currCount >= int(count) {
				break
			}
		}

		isFinished = iterator.Valid()
		return nil
	})
	if err != nil {
		return nstorage.ScanResult{}, nerror.WrapOnly(err)
	}

	var nextKey string
	if len(keys) != 0 {
		nextKey = keys[len(keys)-1]
	}

	return nstorage.ScanResult{
		Finished: isFinished,
		Keys:     keys,
		LastKey:  nextKey,
	}, nil
}

// EachKeyMatch returns all matching results within using giving functions.
func (rd *BadgerStore) EachKeyMatch(regexp string) (keys []string, err error) {
	var generatedRegEx *regexp2.Regexp
	var rgErr error

	if len(regexp) != 0 {
		generatedRegEx, rgErr = regexp2.Compile(regexp)
		if rgErr != nil {
			return nil, nerror.WrapOnly(rgErr)
		}
	}

	keys = make([]string, 0, 10)
	err = rd.Db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			if generatedRegEx != nil && !generatedRegEx.Match(item.Key()) {
				continue
			}
			keys = append(keys, nunsafe.Bytes2String(item.Key()))
		}
		return nil
	})
	return
}

func (rd *BadgerStore) Count() (int64, error) {
	var count int64
	var readErr = rd.Db.View(func(txn *badger.Txn) error {
		return nil
	})
	if readErr != nil {
		return -1, nerror.WrapOnly(readErr)
	}
	return count, nil
}

// Each runs through all elements for giving store, skipping keys
// in Badger who have no data or an empty byte slice.
//
// Each byte slice provided is only valid for the call of
// the function, after which it becomes invalid as it can
// be re-used for efficient memory management, so ensure to copy
// given byte slice yourself within function to protect against
// undefined behaviour.
func (rd *BadgerStore) Each(fn nstorage.EachItem) error {
	return rd.Db.View(func(txn *badger.Txn) error {
		var iterator = txn.NewIterator(rd.iter)
		defer iterator.Close()

		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			var item = iterator.Item()
			if item.IsDeletedOrExpired() {
				continue
			}
			var stop = false
			var err = item.Value(func(value []byte) error {
				if dataErr := fn(value, string(item.Key())); dataErr != nil {
					stop = true
					if nerror.IsAny(dataErr, nstorage.ErrJustStop) {
						return nil
					}
					return dataErr
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
		item, err := txn.Get(nunsafe.String2Bytes(key))
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

// GetAnyKeys returns a list of values for any of the key's found.
// Unless a specific error occurred retrieving the value of a key, if a
// key is not found then it is ignored and a nil is set in it's place.
func (rd *BadgerStore) GetAnyKeys(keys ...string) ([][]byte, error) {
	var values = make([][]byte, 0, len(keys))
	if err := rd.Db.View(func(txn *badger.Txn) error {
		for _, key := range keys {
			item, err := txn.Get(nunsafe.String2Bytes(key))
			if err != nil {
				return nerror.WrapOnly(err)
			}
			if item.IsDeletedOrExpired() {
				values = append(values, nil)
				continue
			}
			var value, verr = item.ValueCopy(nil)
			if verr != nil {
				return nerror.WrapOnly(verr)
			}
			values = append(values, value)
		}

		return nil
	}); err != nil {
		return values, err
	}
	return values, nil
}

// GetAllKeys returns a list of values for any of the key's found.
// if the value of a key is not found then we stop immediately, returning
// an error and the current set of items retreived.
func (rd *BadgerStore) GetAllKeys(keys ...string) ([][]byte, error) {
	var values = make([][]byte, 0, len(keys))
	if err := rd.Db.View(func(txn *badger.Txn) error {
		for _, key := range keys {
			item, err := txn.Get(nunsafe.String2Bytes(key))
			if err != nil {
				return nerror.WrapOnly(err)
			}
			if item.IsDeletedOrExpired() {
				return nerror.New("not found")
			}
			var value, verr = item.ValueCopy(nil)
			if verr != nil {
				return nerror.WrapOnly(verr)
			}
			values = append(values, value)
		}

		return nil
	}); err != nil {
		return values, err
	}
	return values, nil
}

// Get returns giving session stored with giving key, returning an
// error if not found.
func (rd *BadgerStore) Get(key string) ([]byte, error) {
	var value []byte
	if err := rd.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(nunsafe.String2Bytes(key))
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
		op.Key = nunsafe.String2Bytes(key)

		if expiration > 0 {
			op.WithTTL(expiration)
		}

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
		var item, err = txn.Get(nunsafe.String2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		ttl = ttlDur(item.ExpiresAt(), 0)
		if ttl < 0 {
			ttl *= -1
		}
		return nil
	})
	return ttl, err
}

// Close updates to disk.
func (rd *BadgerStore) Close() error {
	return rd.Db.Close()
}

// Sync updates to disk.
func (rd *BadgerStore) Sync() error {
	return rd.Db.Sync()
}

// ExtendTTL resets new TTL for giving key if it has not expired and is still accessible.
//
// A expiration value of zero means to persist the giving key.
func (rd *BadgerStore) ExtendTTL(key string, expiration time.Duration) error {
	return rd.Db.Update(func(txn *badger.Txn) error {
		var item, err = txn.Get(nunsafe.String2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		var lastTTL = ttlDur(item.ExpiresAt(), 0)
		if lastTTL < 0 {
			lastTTL *= -1
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return nerror.Wrap(err, "failed to delete")
		}

		// delete old key
		if err := txn.Delete(nunsafe.String2Bytes(key)); err != nil {
			return nerror.Wrap(err, "failed to delete")
		}

		var op badger.Entry
		op.Value = value
		op.Key = nunsafe.String2Bytes(key)

		if expiration > 0 {
			op.WithTTL(lastTTL + expiration)
		}

		if err := txn.SetEntry(&op); err != nil {
			return nerror.Wrap(err, "failed to save")
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
		var item, err = txn.Get(nunsafe.String2Bytes(key))
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
		op.Key = nunsafe.String2Bytes(key)

		if expiration > 0 {
			op.WithTTL(expiration)
		}

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
		var item, err = txn.Get(nunsafe.String2Bytes(key))
		if err != nil {
			return nerror.Wrap(err, "Failed to retrieve key")
		}
		if item.IsDeletedOrExpired() {
			return nerror.New("not found, possibly expired")
		}

		var op badger.Entry
		op.Value = data
		op.Key = nunsafe.String2Bytes(key)

		if expiration > 0 {
			op.WithTTL(expiration)
		}

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
		var item, err = txn.Get(nunsafe.String2Bytes(key))
		if err != nil {
			return err
		}

		old = make([]byte, item.ValueSize())
		old, err = item.ValueCopy(old)
		if err != nil {
			return err
		}

		return txn.Delete(nunsafe.String2Bytes(key))
	})
	return old, err
}

// RemoveKeys removes all keys found in the store.
func (rd *BadgerStore) RemoveKeys(keys ...string) error {
	if err := rd.Db.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			err := txn.Delete(nunsafe.String2Bytes(key))
			if err != nil && err != badger.ErrEmptyKey && err != badger.ErrKeyNotFound {
				return nerror.WrapOnly(err)
			}
			return err
		}

		return nil
	}); err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

// *****************************************************
// internal methods
// *****************************************************

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
