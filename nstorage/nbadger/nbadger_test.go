package nbadger

import (
	"testing"

	"github.com/dgraph-io/badger/v2"
	"github.com/influx6/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

func TestStoreWithBadgerStoreGetAnykeys(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreGetAnykeys(t, store)
}

func TestStoreWithBadgerStoreGetAllkeys(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreGetAllkeys(t, store)
}

func TestStoreWithBadgerStoreFindEach(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindEach(t, store)
}

func TestStoreWithBadgerStoreEachKeyPrefix(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindPrefix(t, store)
}

func TestStoreWithBadgerStoreFindAll(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindAll(t, store)
}

func TestStoreWithBadgerStore(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestBadgerExpiryStore(t *testing.T) {
	var ops = badger.DefaultOptions("").WithInMemory(true)
	var store, err = NewBadgerStore(ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
