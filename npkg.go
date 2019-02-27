package npkg

import (
	"errors"
	"time"
)

const (
	// DEBUGKey defines the key/flag name used to identify whether
	// debugging mode is on.
	DEBUGKey = "DEBUG"

	// TRACEKey defines the key/flag name used to identify whether
	// tracing is enabled.
	TRACEKey = "TRACE_ENABLED"
)

var (
	// ErrUnencodable is returned when type can not be encoded automatically.
	ErrUnencodable = errors.New("value can not be encodable and does not implement Encodable interface")

	// ErrUndecodable is returned when type can not be encoded automatically.
	ErrUndecodable = errors.New("value can not be decoded to and does not implement DecodableObject/List interface")
)

//************************************************************
// Getters
//************************************************************

// Getter embodies what we expect from setting/configuration provider getters.
type Getter interface {
	Has(name string) bool

	Int(string) int
	Bool(string) bool
	Int8(string) int8
	Int16(string) int16
	Int32(string) int32
	Int64(string) int64
	Float32(string) float32
	Float64(string) float64
	String(name string) string
	Duration(string) time.Duration
	Interface(string) interface{}
}

//************************************************************
// Decoder
//************************************************************

// DecodableObject defines what we expect the a decodable type for object should provide.
type DecodableObject interface {
	DecodeKey(decoder Decoder, k string) error
}

// DecodableList defines what we expect the a decodable type for list should provide.
type DecodableList interface {
	DecodeIndex(decoder Decoder, index int64, total int64) error
}

// Decoder defines an interface for what we expect a object decoder
// to provide for key value pairs.
type Decoder interface {
	Int(v *int) error
	UInt(v *uint) error
	Bool(v *bool) error
	Int8(v *int8) error
	Hex(v *string) error
	UInt8(v *uint8) error
	Int16(v *int16) error
	UInt16(v *uint16) error
	Int32(v *int32) error
	UInt32(v *uint32) error
	Int64(v *int64) error
	UInt64(v *uint64) error
	String(v *string) error
	Float64(v *float64) error
	Float32(v *float32) error
	Base64(v *int64, bs int) error

	List(DecodableList) error
	Object(DecodableObject) error
}

// DecodeList attempts to as best as possible decode giving value in
// decoder to provided interface object of matching type.
// If it fails to, it returns an error.
func DecodeList(dec Decoder, v interface{}) error {
	var err error
	switch vt := v.(type) {
	case *string:
		err = dec.String(vt)
	case **string:
		if vt == nil {
			*vt = new(string)
		}
		err = dec.String(*vt)
	case *int:
		err = dec.Int(vt)
	case **int:
		if vt == nil {
			*vt = new(int)
		}
		err = dec.Int(*vt)
	case *int8:
		err = dec.Int8(vt)
	case **int8:
		if vt == nil {
			*vt = new(int8)
		}
		err = dec.Int8(*vt)
	case *int16:
		err = dec.Int16(vt)
	case **int16:
		if vt == nil {
			*vt = new(int16)
		}
		err = dec.Int16(*vt)
	case *int32:
		err = dec.Int32(vt)
	case **int32:
		if vt == nil {
			*vt = new(int32)
		}
		err = dec.Int32(*vt)
	case *int64:
		err = dec.Int64(vt)
	case **int64:
		if vt == nil {
			*vt = new(int64)
		}
		err = dec.Int64(*vt)
	case *uint8:
		err = dec.UInt8(vt)
	case **uint8:
		if vt == nil {
			*vt = new(uint8)
		}
		err = dec.UInt8(*vt)
	case *uint16:
		err = dec.UInt16(vt)
	case **uint16:
		if vt == nil {
			*vt = new(uint16)
		}
		err = dec.UInt16(*vt)
	case *uint32:
		err = dec.UInt32(vt)
	case **uint32:
		if vt == nil {
			*vt = new(uint32)
		}
		err = dec.UInt32(*vt)
	case *uint64:
		err = dec.UInt64(vt)
	case **uint64:
		if vt == nil {
			*vt = new(uint64)
		}
		err = dec.UInt64(*vt)
	case *float64:
		err = dec.Float64(vt)
	case **float64:
		if vt == nil {
			*vt = new(float64)
		}
		err = dec.Float64(*vt)
	case *float32:
		err = dec.Float32(vt)
	case **float32:
		if vt == nil {
			*vt = new(float32)
		}
		err = dec.Float32(*vt)
	case *bool:
		err = dec.Bool(vt)
	case **bool:
		if vt == nil {
			*vt = new(bool)
		}
		err = dec.Bool(*vt)
	case DecodableObject:
		err = dec.Object(vt)
	case DecodableList:
		err = dec.List(vt)
	default:
		err = ErrUndecodable
	}
	return err
}

//************************************************************
// Encoder
//************************************************************

// EncodableObject defines what we expect the a Encodable type should provide.
type EncodableObject interface {
	EncodeObject(encoder Encoder) error
}

// EncodableList defines what we expect the a Encodable list type should provide.
type EncodableList interface {
	EncodeList(encoder ListEncoder) error
}

// ListEncoder defines an interface which defines methods for items into
// a underline list encoding.
type ListEncoder interface {
	AddInt(v int) error
	AddBool(v bool) error
	AddUInt(v uint) error
	AddInt8(v int8) error
	AddInt16(v int16) error
	AddInt32(v int32) error
	AddInt64(v int64) error
	AddUInt8(v uint8) error
	AddUInt16(v uint16) error
	AddUInt32(v uint32) error
	AddUInt64(v uint64) error
	AddString(v string) error
	AddFloat64(v float64) error
	AddFloat32(v float32) error
	AddBase64(v int64, b int) error

	AddList(list EncodableList) error
	AddObject(object EncodableObject) error
	AddObjectWith(fn func(encoder Encoder) error) error
	AddListWith(fn func(encoder ListEncoder) error) error
}

// Encoder encodes giving key - value pairs as it's implementer
// decides providing a encapsulated formatting of data as it sees fit.
type Encoder interface {
	Int(k string, v int) error
	UInt(k string, v uint) error
	Bool(k string, v bool) error
	Int8(k string, v int8) error
	Hex(k string, v string) error
	UInt8(k string, v uint8) error
	Int16(k string, v int16) error
	UInt16(k string, v uint16) error
	Int32(k string, v int32) error
	UInt32(k string, v uint32) error
	Int64(k string, v int64) error
	UInt64(k string, v uint64) error
	String(k string, v string) error
	Float64(k string, v float64) error
	Float32(k string, v float32) error
	Base64(k string, v int64, b int) error

	List(k string, list EncodableList) error
	Object(k string, object EncodableObject) error
	ObjectFor(k string, fx func(Encoder) error) error
	ListFor(k string, fx func(ListEncoder) error) error
}

// Encodable defines a type which exposes a method to encode it's internals
// with provided encoder.
type Encodable interface {
	Encode(encoder Encoder) error
}

// EncodeKV encodes a giving key-value pair into provided encoder based
func EncodeKV(enc Encoder, k string, v interface{}) error {
	switch vt := v.(type) {
	case EncodableObject:
		return enc.Object(k, vt)
	case EncodableList:
		return enc.List(k, vt)
	case string:
		return enc.String(k, vt)
	case bool:
		return enc.Bool(k, vt)
	case int:
		return enc.Int(k, vt)
	case uint:
		return enc.UInt(k, vt)
	case int64:
		return enc.Int64(k, vt)
	case int32:
		return enc.Int32(k, vt)
	case int16:
		return enc.Int16(k, vt)
	case int8:
		return enc.Int8(k, vt)
	case uint64:
		return enc.UInt64(k, vt)
	case uint32:
		return enc.UInt32(k, vt)
	case uint16:
		return enc.UInt16(k, vt)
	case uint8:
		return enc.UInt8(k, vt)
	case float64:
		return enc.Float64(k, vt)
	case float32:
		return enc.Float32(k, vt)
	default:
		return ErrUnencodable
	}
}

// EncodeList encodes a giving key-value pair into provided encoder based
func EncodeList(enc ListEncoder, v interface{}) error {
	switch vt := v.(type) {
	case EncodableObject:
		return enc.AddObject(vt)
	case EncodableList:
		return enc.AddList(vt)
	case string:
		return enc.AddString(vt)
	case bool:
		return enc.AddBool(vt)
	case int:
		return enc.AddInt(vt)
	case uint:
		return enc.AddUInt(vt)
	case int64:
		return enc.AddInt64(vt)
	case int32:
		return enc.AddInt32(vt)
	case int16:
		return enc.AddInt16(vt)
	case int8:
		return enc.AddInt8(vt)
	case uint64:
		return enc.AddUInt64(vt)
	case uint32:
		return enc.AddUInt32(vt)
	case uint16:
		return enc.AddUInt16(vt)
	case uint8:
		return enc.AddUInt8(vt)
	case float64:
		return enc.AddFloat64(vt)
	case float32:
		return enc.AddFloat32(vt)
	default:
		return ErrUnencodable
	}
}
