package natomic

import (
	"sync"
	"sync/atomic"
)

//******************************************************
// Switches
//******************************************************

// BoolSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a BoolSwitch.
type BoolSwitch struct {
	m  sync.Mutex
	fl bool
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *BoolSwitch) Flip(b bool) {
	f.m.Lock()
	f.fl = b
	f.m.Unlock()
}

// State returns the current state of the switch.
func (f *BoolSwitch) Read() bool {
	var state bool
	f.m.Lock()
	state = f.fl
	f.m.Unlock()
	return state
}

// IntSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a IntSwitch.
type IntSwitch struct {
	fl int64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *IntSwitch) Flip(b int64) {
	atomic.StoreInt64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *IntSwitch) Read() int64 {
	return atomic.LoadInt64(&f.fl)
}

// UintSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a UintSwitch.
type UintSwitch struct {
	fl uint64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *UintSwitch) Flip(b uint64) {
	atomic.StoreUint64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *UintSwitch) Read() uint64 {
	return atomic.LoadUint64(&f.fl)
}

// IntCounter defines a concurrent safe bool switch.
//
// Once created, you must not copy a IntCounter.
type IntCounter struct {
	fl int64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *IntCounter) Add(b int64) {
	atomic.AddInt64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *IntCounter) Read() int64 {
	return atomic.LoadInt64(&f.fl)
}

// UintCounter defines a concurrent safe bool switch.
//
// Once created, you must not copy a UintCounter.
type UintCounter struct {
	fl uint64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *UintCounter) Add(b uint64) {
	atomic.AddUint64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *UintCounter) Read() uint64 {
	return atomic.LoadUint64(&f.fl)
}
