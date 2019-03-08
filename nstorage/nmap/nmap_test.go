package nmap_test

import (
	"math/rand"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"

	"github.com/gokit/npkg/nstorage/nmap"
)

func TestAnyMap(t *testing.T) {
	var any = nmap.NewAnyMap(10)
	require.Equal(t, uint(10), any.Capacity)
	require.NotNil(t, any)

	any.Set("amount", "20")
	require.True(t, any.Has("amount"))
	require.Equal(t, "20", any.Get("amount"))

	any.Set("amount", "40")
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", any.Get("amount"))
}

func TestStringAnyMap(t *testing.T) {
	var any = nmap.NewStringAnyMap(10)
	require.Equal(t, uint(10), any.Capacity)
	require.NotNil(t, any)

	any.Set("amount", "20")
	require.True(t, any.Has("amount"))
	require.Equal(t, "20", any.Get("amount"))

	any.Set("amount", "40")
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", any.Get("amount"))
}

func TestByteMap(t *testing.T) {
	var any = nmap.NewByteMap(10)
	require.Equal(t, uint(10), any.Capacity)
	require.NotNil(t, any)

	any.Set("amount", string2Bytes("20"))
	require.True(t, any.Has("amount"))
	require.Equal(t, "20", bytes2String(any.Get("amount")))

	any.Set("amount", string2Bytes("40"))
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", bytes2String(any.Get("amount")))
}

func TestStringMap(t *testing.T) {
	var any = nmap.NewStringMap(10)
	require.Equal(t, uint(10), any.Capacity)
	require.NotNil(t, any)

	any.Set("amount", "20")
	require.True(t, any.Has("amount"))
	require.Equal(t, "20", any.Get("amount"))

	any.Set("amount", "40")
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", any.Get("amount"))
}

func BenchmarkAnyMap(b *testing.B) {
	b.ReportAllocs()

	b.Run("set", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewAnyMap()
		any.SetMany(func(cache map[interface{}]interface{}) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = randomString()
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewAnyMap()
		any.Set("i", 1)

		for i := 0; i < b.N; i++ {
			any.Get("i")
		}
	})
}

func BenchmarkStringAnyMap(b *testing.B) {
	b.Run("set", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewStringAnyMap()
		any.SetMany(func(cache map[string]interface{}) {
			for i := 0; i < b.N; i++ {
				cache["a"] = i
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewStringAnyMap()
		any.Set("i", 1)

		for i := 0; i < b.N; i++ {
			any.Get("i")
		}
	})
}

func BenchmarkByteMap(b *testing.B) {
	b.ReportAllocs()

	b.Run("set", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewByteMap(100)
		any.SetMany(func(cache map[string][]byte) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = string2Bytes(randomString())
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewByteMap(100)
		any.Set("i", string2Bytes("1"))

		for i := 0; i < b.N; i++ {
			any.Get("i")
		}
	})
}

func BenchmarkStringMap(b *testing.B) {
	b.ReportAllocs()

	b.Run("set", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewStringMap(100)
		any.SetMany(func(cache map[string]string) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = randomString()
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = nmap.NewStringMap()
		any.Set("i", "1")

		for i := 0; i < b.N; i++ {
			any.Get("i")
		}
	})
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
