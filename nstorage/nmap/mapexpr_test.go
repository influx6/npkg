package nmap

import (
	"testing"

	"github.com/gokit/npkg/nstorage/internal/tharness"
	"github.com/stretchr/testify/require"
)

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
