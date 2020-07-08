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

// Decode attempts to as best as possible decode giving value in
// decoder to provided interface object of matching type.
// If it fails to, it returns an error.
func Decode(dec Decoder, v interface{}) error {
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

// Encoder defines what we expect from a encoder of object and list elements .
type Encoder interface {
	ObjectEncoderMethods
	ListEncoderMethods
}

// Encodable defines a type which exposes a method to encode it's internals
// with provided encoder.
type Encodable interface {
	Encode(encoder Encoder) error
}

// EncodableObject defines what we expect the a Encodable type should provide.
type EncodableObject interface {
	EncodeObject(encoder ObjectEncoder)
}

// EncodableList defines what we expect the a Encodable list type should provide.
type EncodableList interface {
	EncodeList(encoder ListEncoder)
}

type Error interface {
	Err() error
}

// ObjectEncoder embodies what is expected from a encoding type
// implementing key-value pair encoding.
type ObjectEncoder interface {
	Error
	ObjectEncoderMethods
}

type ObjectEncoderMethods interface {
	Int(k string, v int)
	UInt(k string, v uint)
	Bool(k string, v bool)
	Int8(k string, v int8)
	Hex(k string, v string)
	UInt8(k string, v uint8)
	Int16(k string, v int16)
	UInt16(k string, v uint16)
	Int32(k string, v int32)
	UInt32(k string, v uint32)
	Int64(k string, v int64)
	UInt64(k string, v uint64)
	String(k string, v string)
	Float64(k string, v float64)
	Float32(k string, v float32)
	Base64(k string, v int64, b int)

	List(k string, list EncodableList)
	Object(k string, object EncodableObject)
	ObjectFor(k string, fx func(ObjectEncoder))
	ListFor(k string, fx func(ListEncoder))
}

// ListEncoder defines an interface which defines methods for items into
// a underline list encoding.
type ListEncoder interface {
	Error
	ListEncoderMethods
}

type ListEncoderMethods interface {
	AddInt(v int)
	AddBool(v bool)
	AddUInt(v uint)
	AddInt8(v int8)
	AddInt16(v int16)
	AddInt32(v int32)
	AddInt64(v int64)
	AddUInt8(v uint8)
	AddUInt16(v uint16)
	AddUInt32(v uint32)
	AddUInt64(v uint64)
	AddString(v string)
	AddFloat64(v float64)
	AddFloat32(v float32)
	AddBase64(v int64, b int)

	AddList(list EncodableList)
	AddObject(object EncodableObject)
	AddObjectWith(fn func(encoder ObjectEncoder))
	AddListWith(fn func(encoder ListEncoder))
}

// EncodeKV encodes a giving key-value pair into provided encoder based
func EncodeKV(enc ObjectEncoder, k string, v interface{}) error {
	switch vt := v.(type) {
	case EncodableObject:
		enc.Object(k, vt)
	case EncodableList:
		enc.List(k, vt)
	case map[string]interface{}:
		enc.Object(k, EncodableMap(vt))
	case string:
		enc.String(k, vt)
	case bool:
		enc.Bool(k, vt)
	case int:
		enc.Int(k, vt)
	case uint:
		enc.UInt(k, vt)
	case int64:
		enc.Int64(k, vt)
	case int32:
		enc.Int32(k, vt)
	case int16:
		enc.Int16(k, vt)
	case int8:
		enc.Int8(k, vt)
	case uint64:
		enc.UInt64(k, vt)
	case uint32:
		enc.UInt32(k, vt)
	case uint16:
		enc.UInt16(k, vt)
	case uint8:
		enc.UInt8(k, vt)
	case float64:
		enc.Float64(k, vt)
	case float32:
		enc.Float32(k, vt)
	}
	return enc.Err()
}

// EncodeList encodes a giving key-value pair into provided encoder based
func EncodeList(enc ListEncoder, v interface{}) error {
	switch vt := v.(type) {
	case EncodableObject:
		enc.AddObject(vt)
	case EncodableList:
		enc.AddList(vt)
	case map[string]interface{}:
		enc.AddObject(EncodableMap(vt))
	case string:
		enc.AddString(vt)
	case bool:
		enc.AddBool(vt)
	case int:
		enc.AddInt(vt)
	case uint:
		enc.AddUInt(vt)
	case int64:
		enc.AddInt64(vt)
	case int32:
		enc.AddInt32(vt)
	case int16:
		enc.AddInt16(vt)
	case int8:
		enc.AddInt8(vt)
	case uint64:
		enc.AddUInt64(vt)
	case uint32:
		enc.AddUInt32(vt)
	case uint16:
		enc.AddUInt16(vt)
	case uint8:
		enc.AddUInt8(vt)
	case float64:
		enc.AddFloat64(vt)
	case float32:
		enc.AddFloat32(vt)
	}
	return enc.Err()
}

// EncodableMapList defines a map type which implements the EncodableList interface.
// It attempts to encode all properties accordingly else returns an error in regard's
// giving failure.
type EncodableMapList []map[string]interface{}

// EncodableMap implements the EncodableList interface.
func (enc EncodableMapList) EncodeList(encoder ListEncoder) {
	for _, value := range enc {
		encoder.AddObject(EncodableMap(value))
	}
}

// EncodableMap defines a map type which implements the EncodableObject interface.
// It attempts to encode all properties accordingly else returns an error in regard's
// giving failure.
type EncodableMap map[string]interface{}

// EncodableMap implements the EncodableObject interface.
func (enc EncodableMap) EncodeObject(encoder ObjectEncoder) {
	for key, value := range enc {
		_ = EncodeKV(encoder, key, value)
	}
}
