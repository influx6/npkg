package nmap

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestAnyMap(t *testing.T) {
	var any = NewAnyMap(10)
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
	var any = NewStringAnyMap(10)
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
	var any = NewByteMap(10)
	require.Equal(t, uint(10), any.Capacity)
	require.NotNil(t, any)

	any.Set("amount", string2Bytes("20"))
	require.True(t, any.Has("amount"))
	require.Equal(t, "20", bytes2String(any.Get("amount")))

	any.Set("amount", string2Bytes("40"))
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", bytes2String(any.Get("amount")))
}

func TestExpiringByteMap(t *testing.T) {
	var any = NewExpiringByteMap()
	require.NotNil(t, any)

	any.Set("amount", string2Bytes("20"), 0)
	any.Set("rewards", string2Bytes("20"), time.Second)
	any.Set("rewards", string2Bytes("200"), time.Second)

	require.True(t, any.Has("amount"))
	require.Equal(t, time.Duration(0), any.TTL("amount"))
	require.Equal(t, "20", bytes2String(any.Get("amount")))

	any.Set("amount", string2Bytes("40"), time.Second)
	require.True(t, any.Has("amount"))
	require.Equal(t, "40", bytes2String(any.Get("amount")))

	var before = any.TTL("amount")
	<-time.After(time.Second)
	require.Empty(t, any.Get("amount"))
	require.NotEqual(t, before, any.TTL("amount"))
}

func TestStringMap(t *testing.T) {
	var any = NewStringMap(10)
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
		var any = NewAnyMap()
		any.SetMany(func(cache map[interface{}]interface{}) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = randomString()
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewAnyMap()
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
		var any = NewStringAnyMap()
		any.SetMany(func(cache map[string]interface{}) {
			for i := 0; i < b.N; i++ {
				cache["a"] = i
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewStringAnyMap()
		any.Set("i", 1)

		for i := 0; i < b.N; i++ {
			any.Get("i")
		}
	})
}

func BenchmarkExpiringByteMap(b *testing.B) {
	b.ReportAllocs()

	b.Run("set", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewExpiringByteMap(100)
		any.SetMany(func(cache map[string]ExpiringValue) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = NewExpiringValue(string2Bytes(randomString()), 0)
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewExpiringByteMap(100)
		any.Set("i", string2Bytes("1"), 0)

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
		var any = NewByteMap(100)
		any.SetMany(func(cache map[string][]byte) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = string2Bytes(randomString())
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewByteMap(100)
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
		var any = NewStringMap(100)
		any.SetMany(func(cache map[string]string) {
			for i := 0; i < b.N; i++ {
				cache[randomString()] = randomString()
			}
		})
	})

	b.Run("set_get", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		var any = NewStringMap()
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
