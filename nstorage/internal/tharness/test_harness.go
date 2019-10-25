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

func TestByteStore(t *testing.T, store nstorage.ByteStore) {
	require.NoError(t, store.Save("day", string2Bytes("wrecker")))
	var val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))

	var exist bool
	exist, err = store.Exists("day")
	require.NoError(t, err)
	require.True(t, exist)

	require.NoError(t, store.Update("day", string2Bytes("tweeter")))

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
	require.Equal(t, "tweeter", bytes2String(val))
}

func TestExpirableStore(t *testing.T, store nstorage.ExpirableStore) {
	require.NoError(t, store.SaveTTL("day", string2Bytes("wrecker"), time.Second))

	var val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "wrecker", bytes2String(val))

	ttl, err := store.TTL("day")
	require.NoError(t, err)
	require.True(t, time.Second <= ttl)

	require.NoError(t, store.ExtendTTL("day", time.Second))

	ttl, err = store.TTL("day")
	require.NoError(t, err)
	require.True(t, (time.Second) < ttl)

	require.NoError(t, store.UpdateTTL("day", string2Bytes("tweeter"), time.Second))

	val, err = store.Get("day")
	require.NoError(t, err)
	require.Equal(t, "tweeter", bytes2String(val))

	ttl, err = store.TTL("day")
	require.NoError(t, err)
	require.True(t, (time.Second*2) < ttl)

	fmt.Printf("BeforeTTL: %+q\n", ttl)

	<-time.After(time.Second * 10)
	val, err = store.Get("day")
	fmt.Printf("Returned data: %+q\n", val)
	require.Error(t, err)
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
