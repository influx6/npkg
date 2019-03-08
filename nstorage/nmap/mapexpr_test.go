package nmap_test

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

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

	tharnest.TestExpirableStore(t, store)
}
