package nmap

import (
	"sync"
	"sync/atomic"
	"time"
)

//**********************************************************************
// AnyMap
//**********************************************************************

// AnyMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type AnyMap struct {
	Capacity uint
	lock     sync.Mutex
	cache    *atomic.Value
}

// NewAnyMap returns a new instance of a AnyMap.
func NewAnyMap(cap ...uint) *AnyMap {
	var sm AnyMap
	if len(cap) != 0 {
		sm.Capacity = cap[0]
	}
	return &sm
}

// Get returns giving value for key.
func (m *AnyMap) Get(k interface{}) (value interface{}) {
	m.GetMany(func(values map[interface{}]interface{}) {
		value = values[k]
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// WARNING: Never modify the map, ever.
func (m *AnyMap) GetMany(fn func(map[interface{}]interface{})) {
	m.init()
	var cached = m.cache.Load().(map[interface{}]interface{})
	fn(cached)
}

// Has returns true/false giving value exits for key.
func (m *AnyMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[interface{}]interface{}) {
		_, exists = values[k]
	})
	return exists
}

// Set adds giving key into underline map.
func (m *AnyMap) Set(k interface{}, value interface{}) {
	m.SetMany(func(values map[interface{}]interface{}) {
		values[k] = value
	})
}

// SetMany adds giving key into underline map.
func (m *AnyMap) SetMany(fn func(map[interface{}]interface{})) {
	m.init()

	var cached = m.cache.Load().(map[interface{}]interface{})
	var copied = CopyInterfaceKeyMap(cached)
	fn(copied)

	m.lock.Lock()
	m.cache.Store(copied)
	m.lock.Unlock()
}

func (m *AnyMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}

	var newValue atomic.Value
	var store = make(map[interface{}]interface{}, m.Capacity)
	newValue.Store(store)
	m.cache = &newValue
}

//**********************************************************************
// StringAnyMap
//**********************************************************************

// StringAnyMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type StringAnyMap struct {
	Capacity uint
	lock     sync.Mutex
	cache    *atomic.Value
}

// NewStringAnyMap returns a new instance of a StringAnyMap.
func NewStringAnyMap(cap ...uint) *StringAnyMap {
	var sm StringAnyMap
	if len(cap) != 0 {
		sm.Capacity = cap[0]
	}
	return &sm
}

// Get returns giving value for key.
func (m *StringAnyMap) Get(k string) (value interface{}) {
	m.GetMany(func(values map[string]interface{}) {
		value = values[k]
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// WARNING: Never modify the map, ever.
func (m *StringAnyMap) GetMany(fn func(map[string]interface{})) {
	m.init()

	var cached = m.cache.Load().(map[string]interface{})
	fn(cached)
}

// Has returns true/false giving value exits for key.
func (m *StringAnyMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[string]interface{}) {
		_, exists = values[k]
	})
	return exists
}

// Set adds giving key into underline map.
func (m *StringAnyMap) Set(k string, value interface{}) {
	m.SetMany(func(values map[string]interface{}) {
		values[k] = value
	})
}

// SetMany adds giving key into underline map.
func (m *StringAnyMap) SetMany(fn func(map[string]interface{})) {
	m.init()

	var cached = m.cache.Load().(map[string]interface{})
	var copied = CopyStringKeyMap(cached)
	fn(copied)

	m.lock.Lock()
	m.cache.Store(copied)
	m.lock.Unlock()
}

func (m *StringAnyMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}
	var store = make(map[string]interface{}, m.Capacity)

	var newValue atomic.Value
	newValue.Store(store)
	m.cache = &newValue
}

//**********************************************************************
// StringMap
//**********************************************************************

// StringMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type StringMap struct {
	Capacity uint
	lock     sync.Mutex
	cache    *atomic.Value
}

// NewStringMap returns a new instance of a StringMap.
func NewStringMap(cap ...uint) *StringMap {
	var sm StringMap
	if len(cap) != 0 {
		sm.Capacity = cap[0]
	}
	return &sm
}

// Has returns true/false giving value exits for key.
func (m *StringMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[string]string) {
		_, exists = values[k]
	})
	return exists
}

// Get returns giving value for key.
func (m *StringMap) Get(k string) (value string) {
	m.GetMany(func(values map[string]string) {
		value = values[k]
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// WARNING: Never modify the map, ever.
func (m *StringMap) GetMany(fn func(map[string]string)) {
	m.init()
	var cached = m.cache.Load().(map[string]string)
	fn(cached)
}

// Set adds giving key into underline map.
func (m *StringMap) Set(k string, value string) {
	m.SetMany(func(values map[string]string) {
		values[k] = value
	})
}

// SetMany adds giving key into underline map.
func (m *StringMap) SetMany(fn func(map[string]string)) {
	m.init()

	var cached = m.cache.Load().(map[string]string)
	var copied = CopyStringMap(cached)
	fn(copied)

	m.lock.Lock()
	m.cache.Store(copied)
	m.lock.Unlock()
}

func (m *StringMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}
	var store = make(map[string]string, m.Capacity)

	var newValue atomic.Value
	newValue.Store(store)
	m.cache = &newValue
}

//**********************************************************************
// ByteMap
//**********************************************************************

// ByteMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type ByteMap struct {
	Capacity uint
	lock     sync.Mutex
	cache    *atomic.Value
}

// NewByteMap returns a new instance of a ByteMap.
func NewByteMap(cap ...uint) *ByteMap {
	var sm ByteMap
	if len(cap) != 0 {
		sm.Capacity = cap[0]
	}
	return &sm
}

// Has returns true/false giving value exits for key.
func (m *ByteMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[string][]byte) {
		_, exists = values[k]
	})
	return exists
}

// Get returns giving value for key.
//
// Get makes a copy of the content of the key
// returning that, which causes a single allocation,
// use GetMany to access the content of the key directly
// without any copy, but ensure to copy the content as
// necessary to avoid corruption of value.
func (m *ByteMap) Get(k string) (value []byte) {
	m.GetMany(func(values map[string][]byte) {
		if nvalue, ok := values[k]; ok {
			var content = make([]byte, len(nvalue))
			copy(content, nvalue)
			value = content
		}
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// Get makes a copy of the content of the key
// returning that, which causes a single allocation,
// use GetMany to access the content of the key directly
// without any copy, but ensure to copy the content as
// necessary to avoid corruption of value.
//
// WARNING: Never modify the map, ever.
func (m *ByteMap) GetMany(fn func(map[string][]byte)) {
	m.init()
	var cached = m.cache.Load().(map[string][]byte)
	fn(cached)
}

// Set adds giving key into underline map.
func (m *ByteMap) Set(k string, value []byte) {
	m.SetMany(func(values map[string][]byte) {
		var content = make([]byte, len(value))
		copy(content, value)
		values[k] = content
	})
}

// SetMany adds giving key into underline map.
func (m *ByteMap) SetMany(fn func(map[string][]byte)) {
	m.init()

	var cached = m.cache.Load().(map[string][]byte)
	var copied = CopyStringBytesMap(cached)
	fn(copied)

	m.lock.Lock()
	m.cache.Store(copied)
	m.lock.Unlock()
}

func (m *ByteMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}
	var store = make(map[string][]byte, m.Capacity)

	var newValue atomic.Value
	newValue.Store(store)
	m.cache = &newValue
}

//**********************************************************************
// ExpiringByteMap
//**********************************************************************
var zeroTime = time.Time{}

// ExpiringValue defines a type which holds a giving byte value
// string, it has if attached a possible expiring value, which would
// make it unaccessible once expired.
type ExpiringValue struct {
	Value []byte
	when  time.Time
}

// Elapsed returns the current duration left for expiring.
//
// A positive number means there is still time and a negative
// number means it has expired. But zero means no expiration.
func (ne *ExpiringValue) Elapsed() time.Duration {
	if ne.when.IsZero() {
		return 0
	}
	var current = time.Now()
	if current.Before(ne.when) {
		return ne.when.Sub(current)
	}
	return current.Sub(ne.when)
}

// Expired returns true/false if giving value is expired.
func (ne *ExpiringValue) Expired() bool {
	if !ne.when.IsZero() {
		var current = time.Now()
		if current.After(ne.when) {
			return true
		}
	}
	return false
}

// NewExpiringValue returns a new instance of a ExpiringValue.
func NewExpiringValue(value []byte, ttl time.Duration) ExpiringValue {
	var exr ExpiringValue
	exr.Value = value
	if ttl > 0 {
		exr.when = time.Now().Add(ttl)
	}
	return exr
}

// ExpiringByteMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type ExpiringByteMap struct {
	Capacity uint
	lock     sync.Mutex
	cache    *atomic.Value
}

// NewExpiringByteMap returns a new instance of a ExpiringByteMap.
func NewExpiringByteMap(cap ...uint) *ExpiringByteMap {
	var sm ExpiringByteMap
	if len(cap) != 0 {
		sm.Capacity = cap[0]
	}
	return &sm
}

// Has returns true/false giving value exits for key.
func (m *ExpiringByteMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[string]ExpiringValue) {
		_, exists = values[k]
	})
	return exists
}

// Get returns giving value for key.
//
// Get makes a copy of the content of the key
// returning that, which causes a single allocation,
// use GetMany to access the content of the key directly
// without any copy, but ensure to copy the content as
// necessary to avoid corruption of value.
func (m *ExpiringByteMap) Get(k string) (value []byte) {
	m.GetMany(func(values map[string]ExpiringValue) {
		if nvalue, ok := values[k]; ok {
			if nvalue.Expired() {
				return
			}

			var content = make([]byte, len(nvalue.Value))
			copy(content, nvalue.Value)
			value = content
		}
	})
	return
}

// TTL returns the current remaining time before giving key expires.
func (m *ExpiringByteMap) TTL(k string) (value time.Duration) {
	m.GetMany(func(values map[string]ExpiringValue) {
		if nvalue, ok := values[k]; ok {
			value = nvalue.Elapsed()
		}
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// Get makes a copy of the content of the key
// returning that, which causes a single allocation,
// use GetMany to access the content of the key directly
// without any copy, but ensure to copy the content as
// necessary to avoid corruption of value.
//
// You are expected to respect the expiry values of a ExpiringValue
// and ignore any that as expired as a cleanup will be done later.
//
// WARNING: Never modify the map, ever.
func (m *ExpiringByteMap) GetMany(fn func(map[string]ExpiringValue)) {
	m.init()
	var cached = m.cache.Load().(map[string]ExpiringValue)
	fn(cached)
}

// Set adds giving key into underline map.
//
// if expiration is zero then giving value expiration will not be reset but left
// as is.
//
// Set automatically cleans up the map of expired keys.
func (m *ExpiringByteMap) Set(k string, value []byte, expire time.Duration) {
	m.SetMany(func(values map[string]ExpiringValue) {
		if nval, ok := values[k]; ok {
			nval.Value = value
			if expire > 0 {
				nval.when = time.Now().Add(expire)
			}
			values[k] = nval
			return
		}
		values[k] = NewExpiringValue(value, expire)
	})
}

// ExtendTTL extends giving key value expiration by provided value.
//
// A expiration value of zero means to persist the giving key.
func (m *ExpiringByteMap) ExtendTTL(k string, expire time.Duration) {
	m.SetMany(func(values map[string]ExpiringValue) {
		if nval, ok := values[k]; ok {
			if expire != 0 {
				if nval.when.IsZero() {
					nval.when = time.Now().Add(expire)
				} else {
					nval.when = nval.when.Add(expire)
				}
			} else {
				nval.when = zeroTime
			}
			values[k] = nval
			return
		}
	})
}

// ResetTTL resets the expiration of giving key, persisting if duration
// provided is zero.
func (m *ExpiringByteMap) ResetTTL(k string, expire time.Duration) {
	m.SetMany(func(values map[string]ExpiringValue) {
		if nval, ok := values[k]; ok {
			if expire != 0 {
				nval.when = time.Now().Add(expire)
			} else {
				nval.when = zeroTime
			}
			values[k] = nval
			return
		}
	})
}

// SetMany adds giving key into underline map.
func (m *ExpiringByteMap) SetMany(fn func(map[string]ExpiringValue)) {
	m.init()

	var cached = m.cache.Load().(map[string]ExpiringValue)
	var copied = CopyExpiringBytesMap(cached)
	fn(copied)

	m.lock.Lock()
	m.cache.Store(copied)
	m.lock.Unlock()
}

func (m *ExpiringByteMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}
	var store = make(map[string]ExpiringValue, m.Capacity)

	var newValue atomic.Value
	newValue.Store(store)
	m.cache = &newValue
}
