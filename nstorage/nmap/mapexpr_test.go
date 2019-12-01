package nmap

import (
	"testing"

	"github.com/influx6/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

func TestStoreWithBadgerStoreFindEach(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStoreFindEach(t, store)
}

func TestStoreWithBadgerStoreFindAll(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStoreFindAll(t, store)
}

func TestStoreWithBadgerStore(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestBadgerExpiryStore(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
