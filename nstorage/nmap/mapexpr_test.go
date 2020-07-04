package nmap

import (
	"testing"

	"github.com/influx6/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

func TestStoreWithNMapStoreFindEach(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStoreFindEach(t, store)
}

func TestStoreWithNMapStoreFindAll(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStoreFindAll(t, store)
}

func TestStoreWithNMapStore(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestNMapExpiryStore(t *testing.T) {
	var store = NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
