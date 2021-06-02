package nstr

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"strings"
)

func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RandomIDWithFormat(size int, format string, splitFn func(b []byte) []interface{}) string {
	var b, err = RandomBytes(size)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(format, splitFn(b)...)
}

func RandomID() string {
	return RandomIDWithFormat(8, "%x-%x-%x-%x", func(b []byte) []interface{} {
		return []interface{}{
			b[0:2], b[2:4], b[4:6], b[6:8],
		}
	})
}

func RandomText() string {
	return RandomIDWithFormat(8, "%x%x%x%x", func(b []byte) []interface{} {
		return []interface{}{
			b[0:2], b[2:4], b[4:6], b[6:8],
		}
	})
}

var alphabets = "abcdefghijklmnoprz"

func RandomAlphabet() string {
	var in = mrand.Intn(len(alphabets))
	return string(alphabets[in])
}

func RandomAlphabets(count int) string {
	var br = make([]string, count)
	for i := 0; i < count; i++ {
		br[i] = RandomAlphabet()
	}
	return strings.Join(br, "")
}
