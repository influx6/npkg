package tharness

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/influx6/npkg/nstorage"
)

func TestByteStoreFindAll(t *testing.T, store nstorage.QueryableByteStore) {
	for i := 0; i < 10; i++ {
		var key = fmt.Sprintf("day-%d", i)
		require.NoError(t, store.Save(key, string2Bytes(fmt.Sprintf("i"))))
	}

	var keys, keyErr = store.Keys()
	require.NoError(t, keyErr)
	require.Len(t, keys, 10)

	var count int
	var err = store.Find(func(val []byte, k string) bool {
		count++
		return true
	})

	require.NoError(t, err)
	require.Equal(t, 10, count)
}

func TestByteStoreFindEach(t *testing.T, store nstorage.QueryableByteStore) {
	for i := 0; i < 10; i++ {
		var key = fmt.Sprintf("day-%d", i)
		require.NoError(t, store.Save(key, string2Bytes(fmt.Sprintf("i"))))
	}

	var keys, keyErr = store.Keys()
	require.NoError(t, keyErr)
	require.Len(t, keys, 10)

	var count int
	var err = store.Find(func(val []byte, k string) bool {
		if count >= 2 {
			return false
		}
		count++
		return true
	})

	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestByteStore(t *testing.T, store nstorage.ByteStore) {
	require.NoError(t, store.Save("day", string2Bytes("wrecker")))
	var val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))

	var exist bool
	exist, err = store.Exists("day")
	require.NoError(t, err)
	require.True(t, exist)

	var newValue = "tweeter"
	require.NoError(t, store.Update("day", string2Bytes(newValue)))

	exist, err = store.Exists("day")
	require.NoError(t, err)
	require.True(t, exist)

	var keyValue, kerr = store.Get("day")
	require.NoError(t, kerr)
	require.Equal(t, newValue, string(keyValue))

	var count int
	store.Each(func(v []byte, k string) bool {
		count++
		return true
	})

	require.Equal(t, 1, count)

	val, err = store.Remove("day")
	require.NoError(t, err)
	require.Equal(t, "tweeter", bytes2String(val))
}

func TestExpirableStore(t *testing.T, store nstorage.ExpirableStore) {
	require.NoError(t, store.SaveTTL("day", string2Bytes("wrecker"), 3*time.Second))

	var val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))

	ttl, err := store.TTL("day")
	require.NoError(t, err)
	require.True(t, ttl > 0)

	require.NoError(t, store.ExtendTTL("day", time.Second))

	ttl, err = store.TTL("day")
	require.NoError(t, err)
	require.True(t, ttl > 0)

	require.NoError(t, store.UpdateTTL("day", string2Bytes("tweeter"), time.Second))

	val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "tweeter", bytes2String(val))

	var ttl2, terr = store.TTL("day")
	require.NoError(t, terr)
	fmt.Printf("TTL: %s -> %s\n", ttl2, ttl)
	require.True(t, ttl2 < ttl)
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
