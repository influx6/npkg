package nbadger

import (
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/influx6/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

func TestStoreWithBadgerStoreFindEach(t *testing.T) {
	defer os.RemoveAll("./tmp")
	require.NoError(t, os.MkdirAll("./tmp", 0777))

	var ops = badger.DefaultOptions
	ops.Dir = "/tmp/badger"
	ops.ValueDir = "/tmp/badger"

	var store, err = NewBadgerStore("", ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindEach(t, store)
}

func TestStoreWithBadgerStoreFindAll(t *testing.T) {
	defer os.RemoveAll("./tmp")
	require.NoError(t, os.MkdirAll("./tmp", 0777))

	var ops = badger.DefaultOptions
	ops.Dir = "/tmp/badger"
	ops.ValueDir = "/tmp/badger"

	var store, err = NewBadgerStore("", ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStoreFindAll(t, store)
}

func TestStoreWithBadgerStore(t *testing.T) {
	defer os.RemoveAll("./tmp")
	require.NoError(t, os.MkdirAll("./tmp", 0777))

	var ops = badger.DefaultOptions
	ops.Dir = "/tmp/badger"
	ops.ValueDir = "/tmp/badger"

	var store, err = NewBadgerStore("", ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestBadgerExpiryStore(t *testing.T) {
	defer os.RemoveAll("./tmp")
	require.NoError(t, os.MkdirAll("./tmp/badger", 0777))

	var ops = badger.DefaultOptions
	ops.Dir = "/tmp/badger"
	ops.ValueDir = "/tmp/badger"

	var store, err = NewBadgerStore("", ops, badger.DefaultIteratorOptions)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
