package nbadger_test

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/gokit/npkg/nstorage/internal/tharness"
	"github.com/gokit/npkg/nstorage/nbadger"
	"github.com/stretchr/testify/require"
)

func TestStoreWithBadgerStore(t *testing.T) {
	var ops badger.Options
	var store, err = nbadger.NewBadgerStore("testing", ops, badger.DefaultIterator)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestsBadgerExpiryStore(t *testing.T) {
	var ops badger.Options
	var store, err = nbadger.NewBadgerStore("testing", ops, badger.DefaultIterator)
	require.NoError(t, err)
	require.NotNil(t, store)

	tharnest.TestExpirableStore(t, store)
}
