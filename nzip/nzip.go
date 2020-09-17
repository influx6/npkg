package nzip

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/influx6/npkg/nerror"
)

const (
	off = 0
	on  = 1
)

// UnzipBool unzips giving byte slice value into a boolean
func UnzipBool(val []byte) (bool, error) {
	if len(val) != 1 {
		return false, nerror.New("byte slice must be of length 1")
	}
	if val[0] == 0 {
		return false, nil
	}
	return true, nil
}

// ZipBool zips giving boolean into giving byte slice, returning appended byte slice.
func ZipBool(b bool, c []byte) ([]byte, error) {
	if b {
		return append(c, on), nil
	}
	return append(c, off), nil
}

// UnzipFloat32 converts a byte slice into an float32.
func UnzipFloat32(val []byte) (float32, error) {
	dl, n := DecodeVarInt32(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}

	return DecodeFloat32(dl), nil
}

// UnzipFloat64 converts a byte slice into an float64.
func UnzipFloat64(val []byte) (float64, error) {
	dl, n := DecodeVarInt64(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}

	return DecodeFloat64(dl), nil
}

// ZipFloat32 converts a float64 into a byte slice.
func ZipFloat32(val float32, c []byte) ([]byte, error) {
	enc := EncodeVarInt32(EncodeFloat32(val))
	return append(c, enc...), nil
}

// ZipFloat64 converts a float64 into a byte slice.
func ZipFloat64(val float64, c []byte) ([]byte, error) {
	enc := EncodeVarInt64(EncodeFloat64(val))
	return append(c, enc...), nil
}

// UnzipUint converts a byte slice into an int.
func UnzipUint(val []byte) (uint, error) {
	var dl, n = DecodeVarInt64(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}
	return uint(dl), nil
}

// UnzipInt converts a byte slice into an int.
func UnzipInt(val []byte) (int, error) {
	var dl, n = DecodeVarInt64(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}
	return int(dl), nil
}

// UnzipInt converts a byte slice into an uint32.
func UnzipUint32(val []byte) (uint32, error) {
	dl, n := DecodeVarInt32(val)
	if n == 0 {
		return 0, nerror.New("not found")
	}
	return uint32(dl), nil
}

// UnzipInt converts a byte slice into an int32.
func UnzipInt32(val []byte) (int32, error) {
	dl, n := DecodeVarInt32(val)
	if n == 0 {
		return 0, nerror.New("not found")
	}
	return int32(dl), nil
}

// UnzipInt converts a byte slice into an int64.
func UnzipInt64(val []byte) (int64, error) {
	dl, n := DecodeVarInt64(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}
	return int64(dl), nil
}

// UnzipUint64 converts a byte slice into an uint64.
func UnzipUint64(val []byte) (uint64, error) {
	dl, n := DecodeVarInt64(val)
	if n == 0 {
		return 0, nerror.New("failed to decode float32")
	}
	return uint64(dl), nil
}

// UnzipUint16 converts a byte slice into an uint64.
func UnzipUint16(val []byte) (uint16, error) {
	dl, err := DecodeUint16FromBytes(val)
	if err != nil {
		return 0, err
	}

	return dl, nil
}

// UnzipInt16 converts a byte slice into an uint64.
func UnzipInt16(val []byte) (int16, error) {
	dl, err := DecodeUint16FromBytes(val)
	if err != nil {
		return 0, err
	}

	return int16(dl), nil
}

// ZipInt converts provided value in the range of {int, uint} types
// in all supported arch of 8,16,32,64 into a byte slice.
func ZipInt(b interface{}, c []byte) ([]byte, error) {
	switch val := b.(type) {
	case uint:
		if val < math.MaxUint32 {
			return append(c, EncodeVarInt32(uint32(val))...), nil
		} else {
			return append(c, EncodeVarInt64(uint64(val))...), nil
		}
	case int:
		if val < math.MaxInt32 {
			return append(c, EncodeVarInt32(uint32(val))...), nil
		} else {
			return append(c, EncodeVarInt64(uint64(val))...), nil
		}
	case uint8:
		return append(c, val), nil
	case uint16:
		return append(c, EncodeUInt16(val)...), nil
	case uint32:
		return append(c, EncodeVarInt32(val)...), nil
	case uint64:
		return append(c, EncodeVarInt64(val)...), nil
	case int8:
		return append(c, uint8(val)), nil
	case int16:
		return append(c, EncodeUInt16(uint16(val))...), nil
	case int32:
		return append(c, EncodeVarInt32(uint32(val))...), nil
	case int64:
		return append(c, EncodeVarInt64(uint64(val))...), nil
	}
	return nil, nerror.New("type is not a range of int/uint types")
}

// UnzipTime converts byte slice into a time.Time object using time.RFC3339 as format.
func UnzipTime(val []byte) (time.Time, error) {
	return UnzipTimeWithFormat(val, time.RFC3339)
}

// UnzipTimeWithFormat converts byte slice into a time.Time object using provided format string.
func UnzipTimeWithFormat(val []byte, format string) (time.Time, error) {
	return time.Parse(format, string(val))
}

// ZipTime converts giving time.Time object into a string using time.RFC3339 format.
func ZipTime(b time.Time, c []byte) ([]byte, error) {
	return ZipTimeWithFormat(time.RFC3339, b, c)
}

// ZipTimeWithFormat converts giving time.Time object into a string using a giving format.
func ZipTimeWithFormat(format string, b time.Time, c []byte) ([]byte, error) {
	formatted := b.Format(format)
	return append(c, formatted...), nil
}

//******************************************
// Codec Functions
//******************************************

// DecodeInt16FromBytes attempts to decode provided byte slice
// into a int16 ensuring that it has minimum length of 2.
// It uses binary.BigEndian.
func DecodeInt16FromBytes(val []byte) (int16, error) {
	de, err := DecodeUint16FromBytes(val)
	return int16(de), err
}

// DecodeUint16FromBytes attempts to decode provided byte slice
// into a uint16 ensuring that it has minimum length of 2.
// It uses binary.BigEndian.
func DecodeUint16FromBytes(val []byte) (uint16, error) {
	if len(val) < 2 {
		return 0, errors.New("byte slice length too small, must be 2")
	}

	var err error
	defer func() {
		if it := recover(); it != nil {
			err = errors.New("failed to decode byte slice with binary.BigEndian")
		}
	}()
	return binary.BigEndian.Uint16(val), err
}

// DecodeInt64FromBytes attempts to decode provided byte slice
// into a int64 ensuring that it has minimum length of 8.
// It uses binary.BigEndian.
func DecodeInt64FromBytes(val []byte) (int64, error) {
	de, err := DecodeUint64FromBytes(val)
	return int64(de), err
}

// DecodeUint64FromBytes attempts to decode provided byte slice
// into a uint64 ensuring that it has minimum length of 8.
// It uses binary.BigEndian.
func DecodeUint64FromBytes(val []byte) (uint64, error) {
	if len(val) < 8 {
		return 0, errors.New("byte slice length too small, must be 8")
	}

	var err error
	defer func() {
		if it := recover(); it != nil {
			err = errors.New("failed to decode byte slice with binary.BigEndian")
		}
	}()
	return binary.BigEndian.Uint64(val), err
}

// DecodeInt32FromBytes attempts to decode provided byte slice
// into a int32 ensuring that it has minimum length of 4.
// It uses binary.BigEndian.
func DecodeInt32FromByte(val []byte) (int32, error) {
	de, err := DecodeUint32FromBytes(val)
	return int32(de), err
}

// DecodeUint32FromBytes attempts to decode provided byte slice
// into a uint32 ensuring that it has minimum length of 4.
// It uses binary.BigEndian.
func DecodeUint32FromBytes(val []byte) (uint32, error) {
	if len(val) < 4 {
		return 0, errors.New("byte slice length too small, must be 4")
	}

	var err error
	defer func() {
		if it := recover(); it != nil {
			err = errors.New("failed to decode byte slice with binary.BigEndian")
		}
	}()
	return binary.BigEndian.Uint32(val), err
}

// EncodeInt32ToBytes encodes provided uint32 into provided
// byte ensuring byte slice has minimum of length 4.
// It uses binary.BigEndian.
func EncodeInt32ToBytes(f int32, val []byte) error {
	return EncodeUint32ToBytes(uint32(f), val)
}

// EncodeUint16ToBytes encodes provided uint16 into provided
// byte ensuring byte slice has minimum of length 2.
// It uses binary.BigEndian.
func EncodeUint16ToBytes(f uint16, val []byte) error {
	if cap(val) < 2 {
		return errors.New("required 8 length for size")
	}

	binary.BigEndian.PutUint16(val, f)
	return nil
}

// EncodeUint32ToBytes encodes provided uint32 into provided
// byte ensuring byte slice has minimum of length 4.
// It uses binary.BigEndian.
func EncodeUint32ToBytes(f uint32, val []byte) error {
	if cap(val) < 4 {
		return errors.New("required 8 length for size")
	}

	binary.BigEndian.PutUint32(val, f)
	return nil
}

// EncodeInt64ToBytes encodes provided uint64 into provided
// byte ensuring byte slice has minimum of length 8.
// It uses binary.BigEndian.
func EncodeInt64ToBytes(f int64, val []byte) error {
	return EncodeUint64ToBytes(uint64(f), val)
}

// EncodeUint64ToBytes encodes provided uint64 into provided
// byte ensuring byte slice has minimum of length 8.
// It uses binary.BigEndian.
func EncodeUint64ToBytes(f uint64, val []byte) error {
	if cap(val) < 8 {
		return errors.New("required 8 length for size")
	}

	binary.BigEndian.PutUint64(val, f)
	return nil
}

// DecodeFloat32 will decode provided uint64 value which should be in
// standard IEEE 754 binary representation, where it bit has been reversed,
// where having it's exponent appears first. It returns the float32 value.
func DecodeFloat32(f uint32) float32 {
	rbit := bits.ReverseBytes32(f)
	return math.Float32frombits(rbit)
}

// EncodeFloat64 will encode provided float value into the standard
// IEEE 754 binary representation and has it's bit reversed, having
// the exponent appearing first.
func EncodeFloat32(f float32) uint32 {
	fbit := math.Float32bits(f)
	return bits.ReverseBytes32(fbit)
}

// DecodeFloat64 will decode provided uint64 value which should be in
// standard IEEE 754 binary representation, where it bit has been reversed,
// where having it's exponent appears first. It returns the float64 value.
func DecodeFloat64(f uint64) float64 {
	rbit := bits.ReverseBytes64(f)
	return math.Float64frombits(rbit)
}

// EncodeFloat64 will encode provided float value into the standard
// IEEE 754 binary representation and has it's bit reversed, having
// the exponent appearing first.
func EncodeFloat64(f float64) uint64 {
	fbit := math.Float64bits(f)
	return bits.ReverseBytes64(fbit)
}

// EncodeVarInt32 encodes uint32 into a byte slice
// using EncodeVarInt64 after turing uint32 into uin64.
func EncodeVarInt32(x uint32) []byte {
	return EncodeVarInt64(uint64(x))
}

// EncodeUInt16 returns the encoded byte slice of a uint16 value.
func EncodeUInt16(x uint16) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, x)
	return data
}

// EncodeVarInt64 returns the varint encoding of x.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum.
func EncodeVarInt64(x uint64) []byte {
	var buf [10]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

// DecodeVarInt32 encodes uint32 into a byte slice
// using EncodeVarInt64 after turing uint32 into uin64.
func DecodeVarInt32(val []byte) (uint32, int) {
	v, d := DecodeVarInt64(val)
	return uint32(v), d
}

// DecodeUInt16 returns the decoded uint16 of provided byte slice which
// must be of length 2.
func DecodeUInt16(d []byte) uint16 {
	return binary.BigEndian.Uint16(d)
}

// DecodeVarInt64 reads a varint-encoded integer from the slice.
// It returns the integer and the number of bytes consumed, or
// zero if there is not enough.
// This is the format for the
// int32, int64, uint32, uint64, bool.
func DecodeVarInt64(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

//*******************************************************************
// Name Zipping : Code taken from https://github.com/linkedin/goavro
//*******************************************************************

// While slices in Go are never constants, we can initialize them once and reuse
// them many times. We define these slices at library load time and reuse them
// when encoding JSON.
var (
	sliceQuote          = []byte("\\\"")
	sliceBackslash      = []byte("\\\\")
	sliceSlash          = []byte("\\/")
	sliceBackspace      = []byte("\\b")
	sliceFormfeed       = []byte("\\f")
	sliceNewline        = []byte("\\n")
	sliceCarriageReturn = []byte("\\r")
	sliceTab            = []byte("\\t")
	sliceUnicode        = []byte("\\u")
)

func ZipText(buf []byte) (interface{}, []byte, error) {
	buflen := len(buf)
	if buflen < 2 {
		return nil, nil, fmt.Errorf("cannot decode textual bytes: %s", io.ErrShortBuffer)
	}
	if buf[0] != '"' {
		return nil, nil, fmt.Errorf("cannot decode textual bytes: expected initial \"; found: %#U", buf[0])
	}
	var newBytes []byte
	var escaped bool
	// Loop through bytes following initial double quote, but note we will
	// return immediately when find unescaped double quote.
	for i := 1; i < buflen; i++ {
		b := buf[i]
		if escaped {
			escaped = false
			if b2, ok := unescapeSpecialJSON(b); ok {
				newBytes = append(newBytes, b2)
				continue
			}
			if b == 'u' {
				// NOTE: Need at least 4 more bytes to read uint16, but subtract
				// 1 because do not want to count the trailing quote and
				// subtract another 1 because already consumed u but have yet to
				// increment i.
				if i > buflen-6 {
					return nil, nil, fmt.Errorf("cannot decode textual bytes: %s", io.ErrShortBuffer)
				}
				// NOTE: Avro bytes represent binary data, and do not
				// necessarily represent text. Therefore, Avro bytes are not
				// encoded in UTF-16. Each \u is followed by 4 hexadecimal
				// digits, the first and second of which must be 0.
				v, err := parseUint64FromHexSlice(buf[i+3 : i+5])
				if err != nil {
					return nil, nil, fmt.Errorf("cannot decode textual bytes: %s", err)
				}
				i += 4 // absorb 4 characters: one 'u' and three of the digits
				newBytes = append(newBytes, byte(v))
				continue
			}
			newBytes = append(newBytes, b)
			continue
		}
		if b == '\\' {
			escaped = true
			continue
		}
		if b == '"' {
			return newBytes, buf[i+1:], nil
		}
		newBytes = append(newBytes, b)
	}
	return nil, nil, fmt.Errorf("cannot decode textual bytes: expected final \"; found: %#U", buf[buflen-1])
}

func UnzipText(buf []byte) (interface{}, []byte, error) {
	buflen := len(buf)
	if buflen < 2 {
		return nil, nil, fmt.Errorf("cannot decode textual string: %s", io.ErrShortBuffer)
	}
	if buf[0] != '"' {
		return nil, nil, fmt.Errorf("cannot decode textual string: expected initial \"; found: %#U", buf[0])
	}
	var newBytes []byte
	var escaped bool
	// Loop through bytes following initial double quote, but note we will
	// return immediately when find unescaped double quote.
	for i := 1; i < buflen; i++ {
		b := buf[i]
		if escaped {
			escaped = false
			if b2, ok := unescapeSpecialJSON(b); ok {
				newBytes = append(newBytes, b2)
				continue
			}
			if b == 'u' {
				// NOTE: Need at least 4 more bytes to read uint16, but subtract
				// 1 because do not want to count the trailing quote and
				// subtract another 1 because already consumed u but have yet to
				// increment i.
				if i > buflen-6 {
					return nil, nil, fmt.Errorf("cannot decode textual string: %s", io.ErrShortBuffer)
				}
				v, err := parseUint64FromHexSlice(buf[i+1 : i+5])
				if err != nil {
					return nil, nil, fmt.Errorf("cannot decode textual string: %s", err)
				}
				i += 4 // absorb 4 characters: one 'u' and three of the digits

				nbl := len(newBytes)
				newBytes = append(newBytes, []byte{0, 0, 0, 0}...) // grow to make room for UTF-8 encoded rune

				r := rune(v)
				if utf16.IsSurrogate(r) {
					i++ // absorb final hexadecimal digit from previous value

					// Expect second half of surrogate pair
					if i > buflen-6 || buf[i] != '\\' || buf[i+1] != 'u' {
						return nil, nil, errors.New("cannot decode textual string: missing second half of surrogate pair")
					}

					v, err = parseUint64FromHexSlice(buf[i+2 : i+6])
					if err != nil {
						return nil, nil, fmt.Errorf("cannot decode textual string: %s", err)
					}
					i += 5 // absorb 5 characters: two for '\u', and 3 of the 4 digits

					// Get code point by combining high and low surrogate bits
					r = utf16.DecodeRune(r, rune(v))
				}

				width := utf8.EncodeRune(newBytes[nbl:], r) // append UTF-8 encoded version of code point
				newBytes = newBytes[:nbl+width]             // trim off excess bytes
				continue
			}
			newBytes = append(newBytes, b)
			continue
		}
		if b == '\\' {
			escaped = true
			continue
		}
		if b == '"' {
			return string(newBytes), buf[i+1:], nil
		}
		newBytes = append(newBytes, b)
	}
	if escaped {
		return nil, nil, fmt.Errorf("cannot decode textual string: %s", io.ErrShortBuffer)
	}
	return nil, nil, fmt.Errorf("cannot decode textual string: expected final \"; found: %x", buf[buflen-1])
}

func QuoteBytes(buf []byte, datum []byte) ([]byte, error) {
	buf = append(buf, '"') // prefix buffer with double quote
	for _, b := range datum {
		if escaped, ok := escapeSpecialJSON(b); ok {
			buf = append(buf, escaped...)
			continue
		}
		if r := rune(b); r < utf8.RuneSelf && unicode.IsPrint(r) {
			buf = append(buf, b)
			continue
		}
		// This Code Point _could_ be encoded as a single byte, however, it's
		// above standard ASCII range (b > 127), therefore must encode using its
		// four-byte hexadecimal equivalent, which will always start with the
		// high byte 00
		buf = appendUnicodeHex(buf, uint16(b))
	}
	return append(buf, '"'), nil // postfix buffer with double quote
}

func QuoteString(buf []byte, datum string) ([]byte, error) {
	buf = append(buf, '"') // prefix buffer with double quote
	for _, r := range datum {
		if escaped, ok := escapeSpecialJSON(byte(r)); ok {
			buf = append(buf, escaped...)
			continue
		}
		if r < utf8.RuneSelf && unicode.IsPrint(r) {
			buf = append(buf, byte(r))
			continue
		}
		// NOTE: Attempt to encode code point as UTF-16 surrogate pair
		r1, r2 := utf16.EncodeRune(r)
		if r1 != unicode.ReplacementChar || r2 != unicode.ReplacementChar {
			// code point does require surrogate pair, and thus two uint16 values
			buf = appendUnicodeHex(buf, uint16(r1))
			buf = appendUnicodeHex(buf, uint16(r2))
			continue
		}
		// Code Point does not require surrogate pair.
		buf = appendUnicodeHex(buf, uint16(r))
	}
	return append(buf, '"'), nil // postfix buffer with double quote
}

func UnescapeUnicodeString(some string) (string, error) {
	if some == "" {
		return "", nil
	}
	buf := []byte(some)
	buflen := len(buf)
	var i int
	var newBytes []byte
	var escaped bool
	// Loop through bytes following initial double quote, but note we will
	// return immediately when find unescaped double quote.
	for i = 0; i < buflen; i++ {
		b := buf[i]
		if escaped {
			escaped = false
			if b == 'u' {
				// NOTE: Need at least 4 more bytes to read uint16, but subtract
				// 1 because do not want to count the trailing quote and
				// subtract another 1 because already consumed u but have yet to
				// increment i.
				if i > buflen-6 {
					return "", fmt.Errorf("cannot replace escaped characters with UTF-8 equivalent: %s", io.ErrShortBuffer)
				}
				v, err := parseUint64FromHexSlice(buf[i+1 : i+5])
				if err != nil {
					return "", fmt.Errorf("cannot replace escaped characters with UTF-8 equivalent: %s", err)
				}
				i += 4 // absorb 4 characters: one 'u' and three of the digits

				nbl := len(newBytes)
				newBytes = append(newBytes, []byte{0, 0, 0, 0}...) // grow to make room for UTF-8 encoded rune

				r := rune(v)
				if utf16.IsSurrogate(r) {
					i++ // absorb final hexadecimal digit from previous value

					// Expect second half of surrogate pair
					if i > buflen-6 || buf[i] != '\\' || buf[i+1] != 'u' {
						return "", errors.New("cannot replace escaped characters with UTF-8 equivalent: missing second half of surrogate pair")
					}

					v, err = parseUint64FromHexSlice(buf[i+2 : i+6])
					if err != nil {
						return "", fmt.Errorf("cannot replace escaped characters with UTF-8 equivalents: %s", err)
					}
					i += 5 // absorb 5 characters: two for '\u', and 3 of the 4 digits

					// Get code point by combining high and low surrogate bits
					r = utf16.DecodeRune(r, rune(v))
				}

				width := utf8.EncodeRune(newBytes[nbl:], r) // append UTF-8 encoded version of code point
				newBytes = newBytes[:nbl+width]             // trim off excess bytes
				continue
			}
			newBytes = append(newBytes, b)
			continue
		}
		if b == '\\' {
			escaped = true
			continue
		}
		newBytes = append(newBytes, b)
	}
	if escaped {
		return "", fmt.Errorf("cannot replace escaped characters with UTF-8 equivalents: %s", io.ErrShortBuffer)
	}
	return string(newBytes), nil
}

func parseUint64FromHexSlice(buf []byte) (uint64, error) {
	var value uint64
	for _, b := range buf {
		diff := uint64(b - '0')
		if diff < 10 {
			value = (value << 4) | diff
			continue
		}
		b10 := b + 10
		diff = uint64(b10 - 'A')
		if diff < 10 {
			return 0, hex.InvalidByteError(b)
		}
		if diff < 16 {
			value = (value << 4) | diff
			continue
		}
		diff = uint64(b10 - 'a')
		if diff < 10 {
			return 0, hex.InvalidByteError(b)
		}
		if diff < 16 {
			value = (value << 4) | diff
			continue
		}
		return 0, hex.InvalidByteError(b)
	}
	return value, nil
}

func unescapeSpecialJSON(b byte) (byte, bool) {
	// NOTE: The following 8 special JSON characters must be escaped:
	switch b {
	case '"', '\\', '/':
		return b, true
	case 'b':
		return '\b', true
	case 'f':
		return '\f', true
	case 'n':
		return '\n', true
	case 'r':
		return '\r', true
	case 't':
		return '\t', true
	}
	return b, false
}

func appendUnicodeHex(buf []byte, v uint16) []byte {
	// Start with '\u' prefix:
	buf = append(buf, sliceUnicode...)
	// And tack on 4 hexadecimal digits:
	buf = append(buf, hexDigits[(v&0xF000)>>12])
	buf = append(buf, hexDigits[(v&0xF00)>>8])
	buf = append(buf, hexDigits[(v&0xF0)>>4])
	buf = append(buf, hexDigits[(v&0xF)])
	return buf
}

const hexDigits = "0123456789ABCDEF"

func escapeSpecialJSON(b byte) ([]byte, bool) {
	// NOTE: The following 8 special JSON characters must be escaped:
	switch b {
	case '"':
		return sliceQuote, true
	case '\\':
		return sliceBackslash, true
	case '/':
		return sliceSlash, true
	case '\b':
		return sliceBackspace, true
	case '\f':
		return sliceFormfeed, true
	case '\n':
		return sliceNewline, true
	case '\r':
		return sliceCarriageReturn, true
	case '\t':
		return sliceTab, true
	}
	return nil, false
}
