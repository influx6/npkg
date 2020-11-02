package njson

import (
	"sync"

	"github.com/influx6/npkg"
)

var (
	writePool = &sync.Pool{
		New: func() interface{} {
			return &WriteStack{}
		},
	}
)

func GetWriteList(w Writer) *WriteStack {
	var newStack, isStack = writePool.Get().(*WriteStack)
	if !isStack {
		newStack = &WriteStack{}
	}
	newStack.writer = w
	return newStack
}

func GetWriteObject(w Writer) *WriteStack {
	var newStack, isStack = writePool.Get().(*WriteStack)
	if !isStack {
		newStack = &WriteStack{}
	}
	newStack.writer = w
	return newStack
}

func ReleaseWriteObject(ll *WriteStack) {
	writePool.Put(ll)
}

type Writer interface {
	Write(*JSON)
}

// WriteStack implements a simple chainable wrapper around the json object which it sends to a Writer
// when New or End is called.
//
// It is not thread-safe to call New or End.
type WriteStack struct {
	writer Writer
	base   *JSON
}

func (l *WriteStack) New() *WriteStack {
	l.End()
	l.base = JSONB()
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

func (l *WriteStack) Info() *WriteStack {
	return l.Level(INFO)
}

func (l *WriteStack) Debug() *WriteStack {
	return l.Level(DEBUG)
}

func (l *WriteStack) Critical() *WriteStack {
	return l.Level(CRITICAL)
}

func (l *WriteStack) Panic() *WriteStack {
	return l.Level(PANIC)
}

func (l *WriteStack) Warn() *WriteStack {
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

func (l *WriteStack) AddList(v npkg.EncodableList) *WriteStack {
	l.base.AddList(v)
	return l
}

func (l *WriteStack) AddObject(v npkg.EncodableObject) *WriteStack {
	l.base.AddObject(v)
	return l
}

func (l *WriteStack) AddObjectWith(fn func(encoder npkg.ObjectEncoder)) *WriteStack {
	l.base.AddObjectWith(fn)
	return l
}

func (l *WriteStack) AddListWith(fn func(encoder npkg.ListEncoder)) *WriteStack {
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

func (l *WriteStack) List(k string, list npkg.EncodableList) *WriteStack {
	l.base.List(k, list)
	return l
}

func (l *WriteStack) Object(k string, object npkg.EncodableObject) *WriteStack {
	l.base.Object(k, object)
	return l
}

func (l *WriteStack) ObjectFor(k string, fx func(npkg.ObjectEncoder)) *WriteStack {
	l.base.ObjectFor(k, fx)
	return l
}

func (l *WriteStack) ListFor(k string, fx func(npkg.ListEncoder)) *WriteStack {
	l.base.ListFor(k, fx)
	return l
}
