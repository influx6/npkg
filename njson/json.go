package njson

import (
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/gokit/npkg"
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

// List requests allocation for a *JSON from the internal pool returning a *JSON
// object for encoding a json list.
func List(inherits ...func(event npkg.Encoder) error) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 1
	event.reset()

	for _, op := range inherits {
		if err := op(event); err != nil {
			panic(err)
		}
	}
	return event
}

// Object requests allocation for a *JSON from the internal pool returning a *JSON
// object for encoding a json object.
func Object(inherits ...func(event npkg.Encoder) error) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0
	event.reset()

	for _, op := range inherits {
		if err := op(event); err != nil {
			panic(err)
		}
	}

	return event
}

// ObjectWithEmbed returns a new Object which will embed all encoded key-value pairs into a object with the `ctx` value
// as key name.
func ObjectWithEmbed(ctx string, hook func(npkg.ObjectEncoder) error, inherits ...func(event npkg.Encoder) error) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0

	event.reset()
	event.onRelease = func(s []byte) []byte {
		newEvent := logEventPool.Get().(*JSON)
		newEvent.l = 0
		newEvent.reset()

		if hook != nil {
			if err := hook(newEvent); err != nil {
				panic(err)
			}
		}

		newEvent.addBytes(ctx, s)
		newEvent.end()

		content := newEvent.content
		newEvent.content = make([]byte, 0, 512)
		newEvent.release()
		return content
	}

	for _, op := range inherits {
		if err := op(event); err != nil {
			panic(err)
		}
	}
	return event
}

// MessageObject requests allocation for a *JSON from the internal pool returning a *JSON
// object for encoding a json object.
func MessageObject(message string, inherits ...func(event npkg.Encoder) error) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0

	event.reset()
	event.addQuotedString("message", message)
	event.endEntry()

	for _, op := range inherits {
		if err := op(event); err != nil {
			panic(err)
		}
	}

	return event
}

// MessageObjectWithEmbed returns a new Object which will embed all encoded key-value pairs into a object with the `ctx` value
// as key name.
func MessageObjectWithEmbed(message string, ctx string, hook func(npkg.ObjectEncoder) error, inherits ...func(npkg.ObjectEncoder) error) *JSON {
	event := logEventPool.Get().(*JSON)
	event.l = 0

	event.reset()
	event.onRelease = func(s []byte) []byte {
		newEvent := logEventPool.Get().(*JSON)
		newEvent.l = 0
		newEvent.reset()

		newEvent.addQuotedString("message", message)
		newEvent.endEntry()

		if hook != nil {
			if err := hook(newEvent); err != nil {
				panic(err)
			}
		}

		newEvent.addBytes(ctx, s)
		newEvent.end()

		content := newEvent.content
		newEvent.content = make([]byte, 0, 512)
		newEvent.release()
		return content
	}

	for _, op := range inherits {
		if err := op(event); err != nil {
			panic(err)
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
// Each JSON iss retrieved from a pool and will panic if after release/write it is used.
type JSON struct {
	l         int8
	r         uint32
	content   []byte
	onRelease func([]byte) []byte
}

func (l *JSON) AddList(list npkg.EncodableList) error {
	return l.AddListWith(list.EncodeList)
}

func (l *JSON) AddObject(object npkg.EncodableObject) error {
	return l.AddObjectWith(object.EncodeObject)
}

func (l *JSON) List(k string, list npkg.EncodableList) error {
	return l.ListFor(k, list.EncodeList)
}

func (l *JSON) Object(k string, object npkg.EncodableObject) error {
	return l.ObjectFor(k, object.EncodeObject)
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

// Release releases the JSON object back into the pool.
func (l *JSON) Release() {
	l.resetContent()
	l.release()
}

// WriteTo makes no attempt like JSON.Message to preserve the byte slice
// data, as it will reuse the byte slice for future writes, it owns it for
// optimization reasons.
//
// It is expected that the writer will adequately copy or write out contents
// of passed in slice before when it's Write method is called.
func (l *JSON) WriteTo(w io.Writer) (int64, error) {
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

	var n, err = w.Write(l.content)
	l.resetContent()
	l.release()
	return int64(n), err
}

// WriteToFN makes no attempt like JSON.Message to preserve the byte slice
// data, as it will reuse the byte slice for future writes, it owns it for
// optimization reasons.
//
// It is expected that the writer will adequately copy or write out contents
// of passed in slice before when it's Write method is called.
func (l *JSON) WriteToFN(fn func([]byte) error) error {
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

	var err = fn(l.content)
	l.resetContent()
	l.release()
	return err
}

// ObjectFor adds a field name with object value.
func (l *JSON) ObjectFor(name string, handler func(event npkg.ObjectEncoder) error) error {
	l.panicIfList()

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 0
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	var err = handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.addBytes(name, newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	return err
}

// ListFor adds a field name with list value.
func (l *JSON) ListFor(name string, handler func(event npkg.ListEncoder) error) error {
	l.panicIfList()

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 1
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	var err = handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.addBytes(name, newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	return err
}

// AddList adds new list object with provided properties from provided function into
// a new json list format. It will panic if you use it for a object format call.
func (l *JSON) AddListWith(handler func(event npkg.ListEncoder) error) error {
	l.panicIfObject()

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 1
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	var err = handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.addBytesListItem(newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	return err
}

// AddObject adds new object with provided properties from provided function into
// a new json list format. It will panic if you use it for a object format call.
func (l *JSON) AddObjectWith(handler func(event npkg.ObjectEncoder) error) error {
	l.panicIfObject()

	newEvent := logEventPool.Get().(*JSON)
	newEvent.l = 0
	newEvent.reset()

	lastLen := len(newEvent.Buf())
	var err = handler(newEvent)
	afterLen := len(newEvent.Buf())

	if afterLen > lastLen {
		total := len(comma) + len(space)
		newEvent.reduce(total)
	}
	newEvent.end()

	l.addBytesListItem(newEvent.Buf())
	l.endEntry()

	newEvent.resetContent()
	newEvent.release()
	return err
}

// AddString adds a string list item into encoding.
func (l *JSON) AddString(value string) error {
	l.panicIfObject()
	l.addQuotedBytesListItem(string2Bytes(value))
	l.endEntry()
	return nil
}

// AddHex adds a hexed string list item into encoding.
func (l *JSON) AddHex(value string) error {
	l.panicIfObject()
	l.addQuotedBytesListItem(string2Bytes(value))
	l.endEntry()
	return nil
}

// AddBool adds a bool value.
func (l *JSON) AddBool(value bool) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		return strconv.AppendBool(content, value)
	})
	l.endEntry()
	return nil
}

// AddInt adds a int value.
func (l *JSON) AddInt(value int) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddInt8 adds a int8 value.
func (l *JSON) AddInt8(value int8) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddInt16 adds a int16 value.
func (l *JSON) AddInt16(value int16) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddInt32 adds a int32 value.
func (l *JSON) AddInt32(value int32) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, int64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddInt64 adds a int64 value.
func (l *JSON) AddInt64(value int64) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, value, 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddUInt adds a int value.
func (l *JSON) AddUInt(value uint) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddUInt8 adds a int8 value.
func (l *JSON) AddUInt8(value uint8) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddUInt16 adds a int16 value.
func (l *JSON) AddUInt16(value uint16) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddUInt32 adds a int32 value.
func (l *JSON) AddUInt32(value uint32) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertUIntToString(content, uint64(value), 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddUInt64 adds a int64 value.
func (l *JSON) AddUInt64(value uint64) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertUIntToString(content, value, 10)
		return content
	})
	l.endEntry()
	return nil
}

// AddBase64 adds a int64 value.
func (l *JSON) AddBase64(value int64, base int) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertIntToString(content, value, base)
		return content
	})
	l.endEntry()
	return nil
}

// AddFloat64 adds a float64 value.
func (l *JSON) AddFloat64(value float64) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertFloatToString(content, value, 32)
		return content
	})
	l.endEntry()
	return nil
}

// AddFloat32 adds a float32 value.
func (l *JSON) AddFloat32(value float32) error {
	l.panicIfObject()
	l.appendItem(func(content []byte) []byte {
		content = convertFloatToString(content, float64(value), 32)
		return content
	})
	l.endEntry()
	return nil
}

// AddBytes adds a bytes value. The byte is expected to be
// valid JSON, no checks are made to ensure this, you can mess up your JSON
// if you do not use this correctly.
func (l *JSON) AddBytes(value []byte) error {
	l.panicIfObject()
	l.addBytesListItem(value)
	l.endEntry()
	return nil
}

// AddQBytes adds a bytes value. The byte is expected to be
// will be wrapped with quotation.
func (l *JSON) AddQBytes(value []byte) error {
	l.panicIfObject()
	l.addQuotedBytesListItem(value)
	l.endEntry()
	return nil
}

// String adds a field name with string value.
func (l *JSON) String(name string, value string) error {
	l.panicIfList()
	l.addQuotedBytes(name, string2Bytes(value))
	l.endEntry()
	return nil
}

// Hex adds a field name with hex converted string value.
func (l *JSON) Hex(name string, value string) error {
	l.panicIfList()
	l.addQuotedBytes(name, string2Bytes(value))
	l.endEntry()
	return nil
}

// Bytes adds a field name with bytes value. The byte is expected to be
// valid JSON, no checks are made to ensure this, you can mess up your JSON
// if you do not use this correctly.
func (l *JSON) Bytes(name string, value []byte) error {
	l.panicIfList()
	l.addBytes(name, value)
	l.endEntry()
	return nil
}

// QBytes adds a field name with bytes value. The byte is expected to be
// will be wrapped with quotation.
func (l *JSON) QBytes(name string, value []byte) error {
	l.panicIfList()
	l.addQuotedBytes(name, value)
	l.endEntry()
	return nil
}

// Bool adds a field name with bool value.
func (l *JSON) Bool(name string, value bool) error {
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
	return nil
}

// Base64 adds a field name with int value formatted to base n.
func (l *JSON) Base64(name string, value int64, base int) error {
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
	return nil
}

// Int adds a field name with int value.
func (l *JSON) Int(name string, value int) error {
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
	return nil
}

// UInt adds a field name with int value.
func (l *JSON) UInt(name string, value uint) error {
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
	return nil
}

// Int8 adds a field name with int8 value.
func (l *JSON) Int8(name string, value int8) error {
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
	return nil
}

// Int16 adds a field name with int16 value.
func (l *JSON) Int16(name string, value int16) error {
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
	return nil
}

// Int32 adds a field name with int32 value.
func (l *JSON) Int32(name string, value int32) error {
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
	return nil
}

// Int64 adds a field name with int64 value.
func (l *JSON) Int64(name string, value int64) error {
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
	return nil
}

// UInt8 adds a field name with uint8 value.
func (l *JSON) UInt8(name string, value uint8) error {
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
	return nil
}

// UInt16 adds a field name with uint16 value.
func (l *JSON) UInt16(name string, value uint16) error {
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
	return nil
}

// UInt32 adds a field name with uint32 value.
func (l *JSON) UInt32(name string, value uint32) error {
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
	return nil
}

// UInt64 adds a field name with uint64 value.
func (l *JSON) UInt64(name string, value uint64) error {
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
	return nil
}

// Float64 adds a field name with float64 value.
func (l *JSON) Float64(name string, value float64) error {
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
	return nil
}

// Float32 adds a field name with float32 value.
func (l *JSON) Float32(name string, value float32) error {
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
	return nil
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

func (l *JSON) addBytes(k string, v []byte) {
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

func (l *JSON) addBytesListItem(v []byte) {
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
