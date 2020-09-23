package plain

import "sync"

// SafeStringMap defines an implementation which is a safe concurrent
// string map value store
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type SafeStringMap struct {
	Capacity uint
	lock     sync.RWMutex
	cache    map[string]string
}

// NewStringMap returns a new instance of a SafeStringMap.
func NewSafeStringMap(capacity uint) *SafeStringMap {
	var sm SafeStringMap
	sm.Capacity = capacity
	return &sm
}

// Has returns true/false giving value exits for key.
func (m *SafeStringMap) Has(k string) bool {
	var exists bool
	m.GetMany(func(values map[string]string) {
		_, exists = values[k]
	})
	return exists
}

// Get returns giving value for key.
func (m *SafeStringMap) Get(k string) (value string) {
	m.GetMany(func(values map[string]string) {
		value = values[k]
	})
	return
}

// GetMany allows retrieval of many keys from underline map.
//
// WARNING: Never modify the map, ever.
func (m *SafeStringMap) GetMany(fn func(map[string]string)) {
	m.init()

	m.lock.RLock()
	fn(m.cache)
	m.lock.RUnlock()
}

// Set adds giving key into underline map.
func (m *SafeStringMap) Set(k string, value string) {
	m.SetMany(func(values map[string]string) {
		values[k] = value
	})
}

// SetMany adds giving key into underline map.
func (m *SafeStringMap) SetMany(fn func(map[string]string)) {
	m.init()

	m.lock.Lock()
	fn(m.cache)
	m.lock.Unlock()
}

func (m *SafeStringMap) init() {
	m.lock.Lock()
	if m.cache != nil {
		m.lock.Unlock()
		return
	}

	defer m.lock.Unlock()
	if m.Capacity == 0 {
		m.Capacity = 10
	}
	m.cache = make(map[string]string, m.Capacity)
}

// StringMap defines an implementation which during initial
// loading stores all key and value pairs.
//
// It provides a safe, concurrently usable implementation with
// blazing read and write speed.
type StringMap struct {
	Capacity uint
	lock     sync.RWMutex
	cache    map[string]string
}

// NewStringMap returns a new instance of a StringMap.
func NewStringMap(capacity uint) *StringMap {
	var sm StringMap
	sm.Capacity = capacity
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

	fn(m.cache)
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
	fn(m.cache)
}

func (m *StringMap) init() {
	if m.cache != nil {
		return
	}

	if m.Capacity == 0 {
		m.Capacity = 10
	}
	m.cache = make(map[string]string, m.Capacity)
}
