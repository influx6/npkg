package nbadger

import (
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/gokit/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

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
