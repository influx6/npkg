package njson

import (
	"sync"

	"github.com/influx6/npkg"
)

var (
	logPool = &sync.Pool{
		New: func() interface{} {
			return &LogStack{}
		},
	}
)

func GetLogList(log Logger) *LogStack {
	var newStack, isStack = logPool.Get().(*LogStack)
	if !isStack {
		newStack = &LogStack{}
	}
	newStack.logger = log
	return newStack
}

func GetLogObject(log Logger) *LogStack {
	var newStack, isStack = logPool.Get().(*LogStack)
	if !isStack {
		newStack = &LogStack{logger: log}
	}
	newStack.logger = log
	return newStack
}

func ReleaseLogStack(ll *LogStack) {
	logPool.Put(ll)
}

type LogLevel int

const (
	DEBUG    LogLevel = 10
	INFO     LogLevel = 20
	WARN     LogLevel = 40
	CRITICAL LogLevel = 60
	PANIC    LogLevel = 79
)

type Logger interface {
	Log(*JSON)
}

// LogStack implements a simple chainable wrapper around the json object which it sends to a Logger
// when New or End is called.
//
// It is not thread-safe to call New or End.
type LogStack struct {
	logger Logger
	base   *JSON
}

func (l *LogStack) New() *LogStack {
	l.End()
	l.base = JSONB()
	return l
}

func (l *LogStack) End() *LogStack {
	if l.base == nil {
		return l
	}
	l.logger.Log(l.base)
	l.base = nil
	return l
}

func (l *LogStack) Info() *LogStack {
	return l.Level(INFO)
}

func (l *LogStack) Debug() *LogStack {
	return l.Level(DEBUG)
}

func (l *LogStack) Critical() *LogStack {
	return l.Level(CRITICAL)
}

func (l *LogStack) Panic() *LogStack {
	return l.Level(PANIC)
}

func (l *LogStack) Warn() *LogStack {
	return l.Level(WARN)
}

func (l *LogStack) Level(level LogLevel) *LogStack {
	l.Int("_level", int(level))
	return l
}

func (l *LogStack) Message(msg string) *LogStack {
	l.String("_message", msg)
	return l
}

func (l *LogStack) AddInt(v int) *LogStack {
	l.base.AddInt(v)
	return l
}

func (l *LogStack) AddBool(v bool) *LogStack {
	l.base.AddBool(v)
	return l
}

func (l *LogStack) AddUInt(v uint) *LogStack {
	l.base.AddUInt(v)
	return l
}

func (l *LogStack) AddInt8(v int8) *LogStack {
	l.base.AddInt8(v)
	return l
}

func (l *LogStack) AddInt16(v int16) *LogStack {
	l.base.AddInt16(v)
	return l
}

func (l *LogStack) AddInt32(v int32) *LogStack {
	l.base.AddInt32(v)
	return l
}

func (l *LogStack) AddByte(v byte) *LogStack {
	l.base.AddByte(v)
	return l
}

func (l *LogStack) AddInt64(v int64) *LogStack {
	l.base.AddInt64(v)
	return l
}

func (l *LogStack) AddUInt8(v uint8) *LogStack {
	l.base.AddUInt8(v)
	return l
}

func (l *LogStack) AddUInt16(v uint16) *LogStack {
	l.base.AddUInt16(v)
	return l
}

func (l *LogStack) AddUInt32(v uint32) *LogStack {
	l.base.AddUInt32(v)
	return l
}

func (l *LogStack) AddUInt64(v uint64) *LogStack {
	l.base.AddUInt64(v)
	return l
}

func (l *LogStack) AddString(v string) *LogStack {
	l.base.AddString(v)
	return l
}

func (l *LogStack) AddFloat64(v float64) *LogStack {
	l.base.AddFloat64(v)
	return l
}

func (l *LogStack) AddFloat32(v float32) *LogStack {
	l.base.AddFloat32(v)
	return l
}

func (l *LogStack) AddBase64(v int64, b int) *LogStack {
	l.base.AddBase64(v, b)
	return l
}

func (l *LogStack) AddMap(v map[string]interface{}) *LogStack {
	l.base.AddMap(v)
	return l
}

func (l *LogStack) AddStringMap(v map[string]string) *LogStack {
	l.base.AddStringMap(v)
	return l
}

func (l *LogStack) AddList(v npkg.EncodableList) *LogStack {
	l.base.AddList(v)
	return l
}

func (l *LogStack) AddObject(v npkg.EncodableObject) *LogStack {
	l.base.AddObject(v)
	return l
}

func (l *LogStack) AddObjectWith(fn func(encoder npkg.ObjectEncoder)) *LogStack {
	l.base.AddObjectWith(fn)
	return l
}

func (l *LogStack) AddListWith(fn func(encoder npkg.ListEncoder)) *LogStack {
	l.base.AddListWith(fn)
	return l
}

func (l *LogStack) Err() error {
	if l.base == nil {
		return nil
	}
	return l.base.Err()
}

func (l *LogStack) Int(k string, v int) *LogStack {
	l.base.Int(k, v)
	return l
}

func (l *LogStack) UInt(k string, v uint) *LogStack {
	l.base.UInt(k, v)
	return l
}

func (l *LogStack) Bool(k string, v bool) *LogStack {
	l.base.Bool(k, v)
	return l
}

func (l *LogStack) Int8(k string, v int8) *LogStack {
	l.base.Int8(k, v)
	return l
}

func (l *LogStack) Hex(k string, v string) *LogStack {
	l.base.Hex(k, v)
	return l
}

func (l *LogStack) UInt8(k string, v uint8) *LogStack {
	l.base.UInt8(k, v)
	return l
}

func (l *LogStack) Int16(k string, v int16) *LogStack {
	l.base.Int16(k, v)
	return l
}

func (l *LogStack) UInt16(k string, v uint16) *LogStack {
	l.base.UInt16(k, v)
	return l
}

func (l *LogStack) Int32(k string, v int32) *LogStack {
	l.base.Int32(k, v)
	return l
}

func (l *LogStack) UInt32(k string, v uint32) *LogStack {
	l.base.UInt32(k, v)
	return l
}

func (l *LogStack) Int64(k string, v int64) *LogStack {
	l.base.Int64(k, v)
	return l
}

func (l *LogStack) UInt64(k string, v uint64) *LogStack {
	l.base.UInt64(k, v)
	return l
}

func (l *LogStack) String(k string, v string) *LogStack {
	l.base.String(k, v)
	return l
}

func (l *LogStack) Bytes(k string, v []byte) *LogStack {
	l.base.Bytes(k, v)
	return l
}

func (l *LogStack) Float64(k string, v float64) *LogStack {
	l.base.Float64(k, v)
	return l
}

func (l *LogStack) Float32(k string, v float32) *LogStack {
	l.base.Float32(k, v)
	return l
}

func (l *LogStack) Base64(k string, v int64, b int) *LogStack {
	l.base.Base64(k, v, b)
	return l
}

func (l *LogStack) Map(k string, v map[string]interface{}) *LogStack {
	l.base.Map(k, v)
	return l
}

func (l *LogStack) StringMap(k string, v map[string]string) *LogStack {
	l.base.StringMap(k, v)
	return l
}

func (l *LogStack) List(k string, list npkg.EncodableList) *LogStack {
	l.base.List(k, list)
	return l
}

func (l *LogStack) Object(k string, object npkg.EncodableObject) *LogStack {
	l.base.Object(k, object)
	return l
}

func (l *LogStack) ObjectFor(k string, fx func(npkg.ObjectEncoder)) *LogStack {
	l.base.ObjectFor(k, fx)
	return l
}

func (l *LogStack) ListFor(k string, fx func(npkg.ListEncoder)) *LogStack {
	l.base.ListFor(k, fx)
	return l
}
