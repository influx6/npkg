package nmap_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nstorage/internal/tharness"
	"github.com/gokit/npkg/nstorage/nmap"
)

func TestStoreWithBadgerStore(t *testing.T) {
	var store = nmap.NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestByteStore(t, store)
}

func TestsBadgerExpiryStore(t *testing.T) {
	var store = nmap.NewExprByteStore(100)
	require.NotNil(t, store)

	tharness.TestExpirableStore(t, store)
}
