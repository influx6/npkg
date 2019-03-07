package nmap

import (
	"sync"
	"sync/atomic"
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
	var copied = copyInterfaceKeyMap(cached)
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
	var copied = copyStringKeyMap(cached)
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
	var copied = copyStringMap(cached)
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
// StringMap
//**********************************************************************

// copyStringBytesMap returns a new copy of a giving string map.
func copyStringBytesMap(src map[string][]byte) map[string][]byte {
	var dest = make(map[string][]byte, len(src))
	for key, value := range src {
		dest[key] = copyBytes(value)
	}
	return dest
}

// copyStringMap returns a new copy of a giving string map.
func copyStringMap(src map[string]string) map[string]string {
	var dest = make(map[string]string, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// copyStringKeyMap returns a new copy of a giving string keyed map.
func copyStringKeyMap(src map[string]interface{}) map[string]interface{} {
	var dest = make(map[string]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// copyInterfaceKeyMap returns a new copy of a giving interface keyed map.
func copyInterfaceKeyMap(src map[interface{}]interface{}) map[interface{}]interface{} {
	var dest = make(map[interface{}]interface{}, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// copyBytes returns a new copy of giving byte slice.
func copyBytes(bu []byte) []byte {
	var cu = make([]byte, len(bu))
	copy(cu, bu)
	return cu
}
