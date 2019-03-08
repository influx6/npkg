package tharness

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nstorage"
)

func TestByteStore(t *testing.T, store nstorage.ByteStore) {
	require.NoError(t, store.Save("day", string2Bytes("wrecker")))
	var val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))

	require.NoError(t, store.Update("day", string2Bytes("tweeter")))

	var exist bool
	exist, err = store.Exists("day")
	require.NoError(t, err)
	require.True(t, exist)

	var count int
	store.Each(func(v []byte, k string) bool {
		count++
		return true
	})

	require.Equal(t, 1, count)

	val, err = store.Remove("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))
}

func TestExpirableStore(t *testing.T, store nstorage.ExpirableStore) {
	require.NoError(t, store.SaveTTL("day", string2Bytes("wrecker"), time.Second*2))
	require.Equal(t, "wrecker", bytes2String(val))

	ttl, err := store.TTL("day")
	require.NoError(t, err)
	require.True(t, time.Second < ttl)

	require.NoError(t, store.ExtendTTL("day", time.Second*3))
	require.True(t, (time.Second*2) < ttl)

	require.NoError(t, store.UpdateTTL("day", string2Bytes("tweeter"), time.Second))
	require.Equal(t, "tweeter", bytes2String(val))
	require.True(t, (time.Second*3) < ttl)
}

func bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

func string2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}

var alphabets = "abcdefghijklmnoprz"

func randomString() string {
	var in = rand.Intn(len(alphabets))
	return string(alphabets[in])
}
