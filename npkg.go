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

type Error interface {
	Err() error
}

// Encoder defines what we expect from a encoder of object and list elements .
type Encoder interface {
	Error
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
	Error(k string, v error)
	Bytes(k string, v []byte)
	Float64(k string, v float64)
	Float32(k string, v float32)
	Base64(k string, v int64, b int)
	Map(k string, v map[string]interface{})
	StringMap(k string, v map[string]string)
	Formatted(k string, format string, v interface{})

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
	AddByte(v byte)
	AddError(v error)
	AddInt64(v int64)
	AddUInt8(v uint8)
	AddUInt16(v uint16)
	AddUInt32(v uint32)
	AddUInt64(v uint64)
	AddString(v string)
	AddFloat64(v float64)
	AddFloat32(v float32)
	AddBase64(v int64, b int)
	AddMap(v map[string]interface{})
	AddStringMap(v map[string]string)
	AddFormatted(format string, v interface{})

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
	case map[string]string:
		enc.StringMap(k, vt)
	case map[string]interface{}:
		enc.Map(k, vt)
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
	case error:
		enc.Error(k, vt)
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
	case map[string]string:
		enc.AddStringMap(vt)
	case map[string]interface{}:
		enc.AddMap(vt)
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
	case error:
		enc.AddError(vt)
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

type EncodedAnyList []interface{}

func (enc EncodedAnyList) EncodeList(encoder ListEncoder) {
	for _, value := range enc {
		_ = EncodeList(encoder, value)
	}
}

type EncodedList []string

func (enc EncodedList) EncodeList(encoder ListEncoder) {
	for _, value := range enc {
		_ = EncodeList(encoder, value)
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

// EncodableMap defines a map type which implements the EncodableObject interface.
// It attempts to encode all properties accordingly else returns an error in regard's
// giving failure.
type EncodableStringMap map[string]string

// EncodableMap implements the EncodableObject interface.
func (enc EncodableStringMap) EncodeObject(encoder ObjectEncoder) {
	for key, value := range enc {
		_ = EncodeKV(encoder, key, value)
	}
}

type Encoded interface{}

type EncodableObjectFunc func() Encoder

type Writer interface {
	Write(Encoded)
}

// WriteStack implements a simple chainable wrapper around the encodable object which it sends to a Writer
// when New or End is called.
//
// It is not thread-safe to call New or End.
type WriteStack struct {
	writer Writer
	maker  EncodableObjectFunc
	base   Encoder
}

func NewWriteStack(fn EncodableObjectFunc, writer Writer) *WriteStack {
	return &WriteStack{
		writer: writer,
		maker:  fn,
	}
}

func (l *WriteStack) SetWriter(w Writer) *WriteStack {
	l.writer = w
	return l
}

func (l *WriteStack) New() *WriteStack {
	l.End()
	l.base = l.maker()
	return l
}

func (l *WriteStack) End() *WriteStack {
	if l.base == nil {
		return l
	}
	l.writer.Write(l.base)
	l.base = nil
	return l
}

func (l *WriteStack) Formatted(k string, format string, v interface{}) *WriteStack {
	l.base.Formatted(k, format, v)
	return l
}

func (l *WriteStack) AddFormatted(format string, v interface{}) *WriteStack {
	l.base.AddFormatted(format, v)
	return l
}

type LogLevel int

const (
	DEBUG    LogLevel = 10
	INFO     LogLevel = 20
	ERROR    LogLevel = 40
	WARN     LogLevel = 60
	CRITICAL LogLevel = 80
	PANIC    LogLevel = 120
)

func (l *WriteStack) LInfo() *WriteStack {
	return l.Level(INFO)
}

func (l *WriteStack) LError() *WriteStack {
	return l.Level(ERROR)
}

func (l *WriteStack) LDebug() *WriteStack {
	return l.Level(DEBUG)
}

func (l *WriteStack) LCritical() *WriteStack {
	return l.Level(CRITICAL)
}

func (l *WriteStack) LPanic() *WriteStack {
	return l.Level(PANIC)
}

func (l *WriteStack) LWarn() *WriteStack {
	return l.Level(WARN)
}

func (l *WriteStack) Level(level LogLevel) *WriteStack {
	l.Int("_level", int(level))
	return l
}

func (l *WriteStack) Message(msg string) *WriteStack {
	l.String("_message", msg)
	return l
}

func (l *WriteStack) AddInt(v int) *WriteStack {
	l.base.AddInt(v)
	return l
}

func (l *WriteStack) AddBool(v bool) *WriteStack {
	l.base.AddBool(v)
	return l
}

func (l *WriteStack) AddUInt(v uint) *WriteStack {
	l.base.AddUInt(v)
	return l
}

func (l *WriteStack) AddInt8(v int8) *WriteStack {
	l.base.AddInt8(v)
	return l
}

func (l *WriteStack) AddInt16(v int16) *WriteStack {
	l.base.AddInt16(v)
	return l
}

func (l *WriteStack) AddInt32(v int32) *WriteStack {
	l.base.AddInt32(v)
	return l
}

func (l *WriteStack) AddByte(v byte) *WriteStack {
	l.base.AddByte(v)
	return l
}

func (l *WriteStack) AddInt64(v int64) *WriteStack {
	l.base.AddInt64(v)
	return l
}

func (l *WriteStack) AddUInt8(v uint8) *WriteStack {
	l.base.AddUInt8(v)
	return l
}

func (l *WriteStack) AddUInt16(v uint16) *WriteStack {
	l.base.AddUInt16(v)
	return l
}

func (l *WriteStack) AddUInt32(v uint32) *WriteStack {
	l.base.AddUInt32(v)
	return l
}

func (l *WriteStack) AddUInt64(v uint64) *WriteStack {
	l.base.AddUInt64(v)
	return l
}

func (l *WriteStack) AddError(v error) *WriteStack {
	l.base.AddError(v)
	return l
}

func (l *WriteStack) AddString(v string) *WriteStack {
	l.base.AddString(v)
	return l
}

func (l *WriteStack) AddFloat64(v float64) *WriteStack {
	l.base.AddFloat64(v)
	return l
}

func (l *WriteStack) AddFloat32(v float32) *WriteStack {
	l.base.AddFloat32(v)
	return l
}

func (l *WriteStack) AddBase64(v int64, b int) *WriteStack {
	l.base.AddBase64(v, b)
	return l
}

func (l *WriteStack) AddMap(v map[string]interface{}) *WriteStack {
	l.base.AddMap(v)
	return l
}

func (l *WriteStack) AddStringMap(v map[string]string) *WriteStack {
	l.base.AddStringMap(v)
	return l
}

func (l *WriteStack) AddList(v EncodableList) *WriteStack {
	l.base.AddList(v)
	return l
}

func (l *WriteStack) AddObject(v EncodableObject) *WriteStack {
	l.base.AddObject(v)
	return l
}

func (l *WriteStack) AddObjectWith(fn func(encoder ObjectEncoder)) *WriteStack {
	l.base.AddObjectWith(fn)
	return l
}

func (l *WriteStack) AddListWith(fn func(encoder ListEncoder)) *WriteStack {
	l.base.AddListWith(fn)
	return l
}

func (l *WriteStack) Err() error {
	if l.base == nil {
		return nil
	}
	return l.base.Err()
}

func (l *WriteStack) Int(k string, v int) *WriteStack {
	l.base.Int(k, v)
	return l
}

func (l *WriteStack) UInt(k string, v uint) *WriteStack {
	l.base.UInt(k, v)
	return l
}

func (l *WriteStack) Bool(k string, v bool) *WriteStack {
	l.base.Bool(k, v)
	return l
}

func (l *WriteStack) Int8(k string, v int8) *WriteStack {
	l.base.Int8(k, v)
	return l
}

func (l *WriteStack) Hex(k string, v string) *WriteStack {
	l.base.Hex(k, v)
	return l
}

func (l *WriteStack) UInt8(k string, v uint8) *WriteStack {
	l.base.UInt8(k, v)
	return l
}

func (l *WriteStack) Int16(k string, v int16) *WriteStack {
	l.base.Int16(k, v)
	return l
}

func (l *WriteStack) UInt16(k string, v uint16) *WriteStack {
	l.base.UInt16(k, v)
	return l
}

func (l *WriteStack) Int32(k string, v int32) *WriteStack {
	l.base.Int32(k, v)
	return l
}

func (l *WriteStack) UInt32(k string, v uint32) *WriteStack {
	l.base.UInt32(k, v)
	return l
}

func (l *WriteStack) Int64(k string, v int64) *WriteStack {
	l.base.Int64(k, v)
	return l
}

func (l *WriteStack) UInt64(k string, v uint64) *WriteStack {
	l.base.UInt64(k, v)
	return l
}

func (l *WriteStack) String(k string, v string) *WriteStack {
	l.base.String(k, v)
	return l
}

func (l *WriteStack) Error(k string, v error) *WriteStack {
	l.base.Error(k, v)
	return l
}

func (l *WriteStack) Bytes(k string, v []byte) *WriteStack {
	l.base.Bytes(k, v)
	return l
}

func (l *WriteStack) Float64(k string, v float64) *WriteStack {
	l.base.Float64(k, v)
	return l
}

func (l *WriteStack) Float32(k string, v float32) *WriteStack {
	l.base.Float32(k, v)
	return l
}

func (l *WriteStack) Base64(k string, v int64, b int) *WriteStack {
	l.base.Base64(k, v, b)
	return l
}

func (l *WriteStack) Map(k string, v map[string]interface{}) *WriteStack {
	l.base.Map(k, v)
	return l
}

func (l *WriteStack) StringMap(k string, v map[string]string) *WriteStack {
	l.base.StringMap(k, v)
	return l
}

func (l *WriteStack) List(k string, list EncodableList) *WriteStack {
	l.base.List(k, list)
	return l
}

func (l *WriteStack) Object(k string, object EncodableObject) *WriteStack {
	l.base.Object(k, object)
	return l
}

func (l *WriteStack) ObjectFor(k string, fx func(ObjectEncoder)) *WriteStack {
	l.base.ObjectFor(k, fx)
	return l
}

func (l *WriteStack) ListFor(k string, fx func(ListEncoder)) *WriteStack {
	l.base.ListFor(k, fx)
	return l
}
