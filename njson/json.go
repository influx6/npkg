package njson

import (
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/influx6/npkg"
)

var (
	comma            = []byte(",")
	colon            = []byte(":")
	space            = []byte(" ")
	openBlock        = []byte("{")
	closingBlock     = []byte("}")
	openBlockList    = []byte("[")
	closingBlockList = []byte("]")
	doubleQuote      = []byte("\"")
	logEventPool     = sync.Pool{
		New: func() interface{} {
			return &JSON{content: make([]byte, 0, 512), r: 1}
		},
	}
)

var _ npkg.Error = (*JSON)(nil)
var _ npkg.ObjectEncoder = (*JSON)(nil)
var _ npkg.ListEncoder = (*JSON)(nil)

// JSONL creates a json list.
func JSONL(inherits ...func(event npkg.Encoder)) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 1
	event.reset()

	for _, op := range inherits {
		op(event)
		if event.err != nil {
			return event
		}
	}
	return event
}

// JSONB creates a json hash.
func JSONB(inherits ...func(event npkg.Encoder)) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0
	event.reset()

	for _, op := range inherits {
		op(event)
		if event.err != nil {
			return event
		}
	}
	return event
}

// MJSON creates a json object with a message field with provided message.
func MJSON(message string, inherits ...func(event npkg.Encoder)) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0

	event.reset()
	event.addQuotedString("message", message)
	event.endEntry()

	for _, op := range inherits {
		op(event)
		if event.err != nil {
			return event
		}
	}
	return event
}

//************************************************************
// JSON
//************************************************************

var (
	_ npkg.Encoder = (*JSON)(nil)
)

// JSON implements a efficient zero or near zero-allocation as much as possible,
// using a underline non-strict json format to transform log key-value pairs into
// a LogMessage.
//
// Each JSON iss retrieved from a logPool and will panic if after release/write it is used.
type JSON struct {
	err       error
	l         int8
	r         uint32
	content   []byte
	onRelease func([]byte) []byte
}

func (l *JSON) AddStringMap(m map[string]string) {
	l.AddObjectWith(func(event npkg.ObjectEncoder) {
		npkg.EncodableStringMap(m).EncodeObject(event)
	})
}

func (l *JSON) AddMap(m map[string]interface{}) {
	l.AddObjectWith(func(event npkg.ObjectEncoder) {
		npkg.EncodableMap(m).EncodeObject(event)
	})
}

func (l *JSON) StringMap(key string, m map[string]string) {
	l.ObjectFor(key, func(event npkg.ObjectEncoder) {
		npkg.EncodableStringMap(m).EncodeObject(event)
	})
}

func (l *JSON) Map(k string, m map[string]interface{}) {
	l.ObjectFor(k, func(event npkg.ObjectEncoder) {
		npkg.EncodableMap(m).EncodeObject(event)
	})
}

func (l *JSON) AddList(list npkg.EncodableList) {
	l.AddListWith(list.EncodeList)
}

func (l *JSON) AddObject(object npkg.EncodableObject) {
	l.AddObjectWith(object.EncodeObject)
}

func (l *JSON) List(k string, list npkg.EncodableList) {
	l.ListFor(k, list.EncodeList)
}

func (l *JSON) Object(k string, object npkg.EncodableObject) {
	l.ObjectFor(k, object.EncodeObject)
}

func (l *JSON) Err() error {
	return l.err
}

// Message returns the generated JSON of giving *JSON.
func (l *JSON) Message() string {
	if l.released() {
		panic("Re-using released *JSON")
	}

	// remove last comma and space
	total := len(comma) + len(space)
	l.reduce(total)
	l.end()

	if l.onRelease != nil {
		l.content = l.onRelease(l.content)
		l.onRelease = nil
	}

	cn := make([]byte, len(l.content))
	copy(cn, l.content)

	l.resetContent()
	l.release()
	return bytes2String(cn)
}

// Release releases the JSON object back into the logPool.
func (l *JSON) Release() {
	l.resetContent()
	l.release()
}

// WriteTo implements io.WriterTo interface.
func (l *JSON) WriteTo(w io.Writer) (int64, error) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	// if there is an error then talk about it.
	if l.err != nil {
		return -1, l.err
	}

	// remove last comma and space
	total := len(comma) + len(space)
	l.reduce(total)
	l.end()

	if l.onRelease != nil {
		l.content = l.onRelease(l.content)
		l.onRelease = nil
	}

	var n, err = w.Write(l.content)
	l.err = err
	l.resetContent()
	l.release()
	return int64(n), err
}

// ObjectFor adds a field name with object value.
func (l *JSON) ObjectFor(name string, handler func(event npkg.ObjectEncoder)) {
	l.panicIfList()

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 0
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.appendBytes(name, newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	if newEvent.err != nil {
		l.err = newEvent.err
	}
}

// ListFor adds a field name with list value.
func (l *JSON) ListFor(name string, handler func(event npkg.ListEncoder)) {
	l.panicIfList()

	// stop if error
	if l.err != nil {
		return
	}

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 1
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.appendBytes(name, newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	if newEvent.err != nil {
		l.err = newEvent.err
	}
}

// AddList adds new list object with provided properties from provided function into
// a new json list format. It will panic if you use it for a object format call.
func (l *JSON) AddListWith(handler func(event npkg.ListEncoder)) {
	l.panicIfObject()
	// stop if error
	if l.err != nil {
		return
	}

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 1
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.appendBytesList(newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()

	if newEvent.err != nil {
		l.err = newEvent.err
	}
}

// AddObject adds new object with provided properties from provided function into
// a new json list format. It will panic if you use it for a object format call.
func (l *JSON) AddObjectWith(handler func(event npkg.ObjectEncoder)) {
	l.panicIfObject()
	// stop if error
	if l.err != nil {
		return
	}

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 0
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.appendBytesList(newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()

	if newEvent.err != nil {
		l.err = newEvent.err
	}
}

// AddString adds a string list item into encoding.
func (l *JSON) AddString(value string) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.addQuotedBytesListItem(string2Bytes(value))
	l.endEntry()
}

// AddHex adds a hexed string list item into encoding.
func (l *JSON) AddHex(value string) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.addQuotedBytesListItem(string2Bytes(value))
	l.endEntry()
}

// AddBool adds a bool value.
func (l *JSON) AddBool(value bool) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return strconv.AppendBool(content, value)
	})
	l.endEntry()
}

// AddInt adds a int value.
func (l *JSON) AddInt(value int) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, int64(value), 10)
	})
	l.endEntry()
}

// AddInt8 adds a int8 value.
func (l *JSON) AddInt8(value int8) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, int64(value), 10)
	})
	l.endEntry()
}

// AddInt16 adds a int16 value.
func (l *JSON) AddInt16(value int16) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, int64(value), 10)
	})
	l.endEntry()
}

// AddInt32 adds a int32 value.
func (l *JSON) AddInt32(value int32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, int64(value), 10)
	})
	l.endEntry()
}

// AddInt64 adds a int64 value.
func (l *JSON) AddInt64(value int64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, value, 10)
	})
	l.endEntry()
}

// AddUInt adds a int value.
func (l *JSON) AddUInt(value uint) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertUIntToString(content, uint64(value), 10)
	})
	l.endEntry()
}

// AddUInt8 adds a int8 value.
func (l *JSON) AddUInt8(value uint8) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertUIntToString(content, uint64(value), 10)
	})
	l.endEntry()
}

// AddUInt16 adds a int16 value.
func (l *JSON) AddUInt16(value uint16) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertUIntToString(content, uint64(value), 10)
	})
	l.endEntry()
}

// AddUInt32 adds a int32 value.
func (l *JSON) AddUInt32(value uint32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertUIntToString(content, uint64(value), 10)
	})
	l.endEntry()
}

// AddUInt64 adds a int64 value.
func (l *JSON) AddUInt64(value uint64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertUIntToString(content, value, 10)
	})
	l.endEntry()
}

// AddBase64 adds a int64 value.
func (l *JSON) AddBase64(value int64, base int) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertIntToString(content, value, base)
	})
	l.endEntry()
}

// AddFloat64 adds a float64 value.
func (l *JSON) AddFloat64(value float64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertFloatToString(content, value, 32)
	})
	l.endEntry()
}

// AddFloat32 adds a float32 value.
func (l *JSON) AddFloat32(value float32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return convertFloatToString(content, float64(value), 32)
	})
	l.endEntry()
}

// AddBytes adds a bytes value. The byte is expected to be
// valid JSON, no checks are made to ensure this, you can mess up your JSON
// if you do not use this correctly.
func (l *JSON) AddBytes(value []byte) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.addListBytes(value)
	l.endEntry()
}

func (l *JSON) AppendBytes(value []byte) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.appendBytesList(value)
	l.endEntry()
}

// AddQBytes adds a bytes value. The byte is expected to be
// will be wrapped with quotation.
func (l *JSON) AddQBytes(value []byte) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfObject()
	l.addQuotedBytesListItem(value)
	l.endEntry()
}

// String adds a field name with string value.
func (l *JSON) String(name string, value string) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.addQuotedBytes(name, string2Bytes(value))
	l.endEntry()

}

// Hex adds a field name with hex converted string value.
func (l *JSON) Hex(name string, value string) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.addQuotedBytes(name, string2Bytes(value))
	l.endEntry()

}

// Bytes adds a field name with bytes value. The byte is expected to be
// valid JSON, no checks are made to ensure this, you can mess up your JSON
// if you do not use this correctly.
func (l *JSON) Bytes(name string, value []byte) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.addListBytesKV(name, value)
	l.endEntry()
}

// QBytes adds a field name with bytes value. The byte is expected to be
// will be wrapped with quotation.
func (l *JSON) QBytes(name string, value []byte) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.addQuotedBytes(name, value)
	l.endEntry()

}

// Bool adds a field name with bool value.
func (l *JSON) Bool(name string, value bool) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = strconv.AppendBool(content, value)
		return content
	})
	l.endEntry()

}

// Base64 adds a field name with int value formatted to base n.
func (l *JSON) Base64(name string, value int64, base int) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, value, base)
		return content
	})
	l.endEntry()

}

// Int adds a field name with int value.
func (l *JSON) Int(name string, value int) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()

}

// UInt adds a field name with int value.
func (l *JSON) UInt(name string, value uint) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()

}

// Int8 adds a field name with int8 value.
func (l *JSON) Int8(name string, value int8) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()

}

// Int16 adds a field name with int16 value.
func (l *JSON) Int16(name string, value int16) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()

}

// Int32 adds a field name with int32 value.
func (l *JSON) Int32(name string, value int32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()

}

// Int64 adds a field name with int64 value.
func (l *JSON) Int64(name string, value int64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()
}

// UInt8 adds a field name with uint8 value.
func (l *JSON) UInt8(name string, value uint8) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()

}

// UInt16 adds a field name with uint16 value.
func (l *JSON) UInt16(name string, value uint16) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()

}

func (l *JSON) Byte(name string, value byte) {
	l.panicIfList()

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, value)
		return content
	})
	l.endEntry()
}

func (l *JSON) AddByte(value byte) {
	l.panicIfList()

	l.appendItem(func(content []byte) []byte {
		content = append(content, value)
		return content
	})
	l.endEntry()
}

// UInt32 adds a field name with uint32 value.
func (l *JSON) UInt32(name string, value uint32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()

}

// UInt64 adds a field name with uint64 value.
func (l *JSON) UInt64(name string, value uint64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertUIntToString(content, value, 10)
		return content
	})
	l.endEntry()

}

// Float64 adds a field name with float64 value.
func (l *JSON) Float64(name string, value float64) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertFloatToString(content, value, 64)
		return content
	})
	l.endEntry()

}

// Float32 adds a field name with float32 value.
func (l *JSON) Float32(name string, value float32) {
	// stop if error
	if l.err != nil {
		return
	}

	l.panicIfList()
	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, name...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = convertFloatToString(content, float64(value), 32)
		return content
	})
	l.endEntry()

}

// Buf returns the current content of the *JSON.
func (l *JSON) Buf() []byte {
	return l.content
}

func (l *JSON) reset() {
	atomic.StoreUint32(&l.r, 1)
	l.begin()
}

func (l *JSON) reduce(d int) {
	available := len(l.content)
	rem := available - d
	if rem < 0 {
		rem = 0
	}
	l.content = l.content[:rem]
}

func (l *JSON) resetContent() {
	l.content = l.content[:0]
}

func (l *JSON) released() bool {
	return atomic.LoadUint32(&l.r) == 0
}

func (l *JSON) release() {
	atomic.StoreUint32(&l.r, 0)
	logEventPool.Put(l)
}

func (l *JSON) begin() {
	l.appendItem(func(content []byte) []byte {
		if l.l == 1 {
			content = append(content, openBlockList...)
			return content
		}
		content = append(content, openBlock...)
		return content
	})
}

func (l *JSON) panicIfObject() {
	if l.l == 0 {
		panic("unable to use for a json object format")
	}
}

func (l *JSON) panicIfList() {
	if l.l == 1 {
		panic("unable to use for a json list format")
	}
}

func (l *JSON) addQuotedString(k string, v string) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, k...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, doubleQuote...)
		content = append(content, v...)
		content = append(content, doubleQuote...)
		return content
	})
}

func (l *JSON) addQuotedStringListItem(k string, v string) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, v...)
		content = append(content, doubleQuote...)
		return content
	})
}

func (l *JSON) addString(k string, v string) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, k...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, v...)
		return content
	})
}

func (l *JSON) addStringListItem(v string) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, v...)
		return content
	})
}

func (l *JSON) addListBytes(v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, '[')
		content = append(content, v...)
		content = append(content, ']')
		return content
	})
}

func (l *JSON) addListBytesKV(k string, v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, k...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, '[')
		content = append(content, v...)
		content = append(content, ']')
		return content
	})
}

func (l *JSON) addQuotedBytes(k string, v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, k...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, doubleQuote...)
		content = append(content, v...)
		content = append(content, doubleQuote...)
		return content
	})
}

func (l *JSON) addQuotedBytesListItem(v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, v...)
		content = append(content, doubleQuote...)
		return content
	})
}

func (l *JSON) appendBytes(k string, v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, doubleQuote...)
		content = append(content, k...)
		content = append(content, doubleQuote...)
		content = append(content, colon...)
		content = append(content, space...)
		content = append(content, v...)
		return content
	})
}

func (l *JSON) appendBytesList(v []byte) {
	if l.released() {
		panic("Re-using released *JSON")
	}

	l.appendItem(func(content []byte) []byte {
		content = append(content, v...)
		return content
	})
}

func (l *JSON) endEntry() {
	l.appendItem(func(content []byte) []byte {
		content = append(content, comma...)
		content = append(content, space...)
		return content
	})
}

func (l *JSON) end() {
	if l.l == 1 {
		l.appendItem(func(content []byte) []byte {
			return append(content, closingBlockList...)
		})
		return
	}

	l.appendItem(func(content []byte) []byte {
		return append(content, closingBlock...)
	})
}

func (l *JSON) appendItem(cb func([]byte) []byte) {
	l.content = cb(l.content)
}

func convertIntToString(in []byte, v int64, b int) []byte {
	if v == 0 {
		return append(in, "0"...)
	}

	return strconv.AppendInt(in, v, b)
}

func convertFloatToString(in []byte, v float64, pre int) []byte {
	if v == 0 {
		return append(in, "0"...)
	}

	return strconv.AppendFloat(in, v, 'E', -1, pre)
}

func convertUIntToString(in []byte, v uint64, b int) []byte {
	if v == 0 {
		return append(in, "0"...)
	}

	return strconv.AppendUint(in, v, b)
}

//*****************************************************
// unsafe methods
//*****************************************************

func bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

func string2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}
