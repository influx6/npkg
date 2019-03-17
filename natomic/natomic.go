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

// Int64Switch defines a concurrent safe bool switch.
//
// Once created, you must not copy a Int64Switch.
type Int64Switch struct {
	fl int64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *Int64Switch) Flip(b int64) {
	atomic.StoreInt64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *Int64Switch) Read() int64 {
	return atomic.LoadInt64(&f.fl)
}

// Uint64Switch defines a concurrent safe bool switch.
//
// Once created, you must not copy a Uint64Switch.
type Uint64Switch struct {
	fl uint64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *Uint64Switch) Flip(b uint64) {
	atomic.StoreUint64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *Uint64Switch) Read() uint64 {
	return atomic.LoadUint64(&f.fl)
}

// Int64Counter defines a concurrent safe bool switch.
//
// Once created, you must not copy a Int64Counter.
type Int64Counter struct {
	fl int64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *Int64Counter) Add(b int64) {
	atomic.AddInt64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *Int64Counter) Read() int64 {
	return atomic.LoadInt64(&f.fl)
}

// Uint64Counter defines a concurrent safe bool switch.
//
// Once created, you must not copy a Uint64Counter.
type Uint64Counter struct {
	fl uint64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *Uint64Counter) Add(b uint64) {
	atomic.AddUint64(&f.fl, b)
}

// State returns the current state of the switch.
func (f *Uint64Counter) Read() uint64 {
	return atomic.LoadUint64(&f.fl)
}

// IntSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a IntSwitch.
type IntSwitch struct {
	fl int64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *IntSwitch) Flip(b int) {
	atomic.StoreInt64(&f.fl, int64(b))
}

// State returns the current state of the switch.
func (f *IntSwitch) Read() int {
	return int(atomic.LoadInt64(&f.fl))
}

// UintSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a UintSwitch.
type UintSwitch struct {
	fl uint64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *UintSwitch) Flip(b uint) {
	atomic.StoreUint64(&f.fl, uint64(b))
}

// State returns the current state of the switch.
func (f *UintSwitch) Read() uint {
	return uint(atomic.LoadUint64(&f.fl))
}

// IntCounter defines a concurrent safe bool switch.
//
// Once created, you must not copy a IntCounter.
type IntCounter struct {
	fl int64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *IntCounter) Add(b int) {
	atomic.AddInt64(&f.fl, int64(b))
}

// State returns the current state of the switch.
func (f *IntCounter) Read() int {
	return int(atomic.LoadInt64(&f.fl))
}

// UintCounter defines a concurrent safe bool switch.
//
// Once created, you must not copy a UintCounter.
type UintCounter struct {
	fl uint64
}

// Add increments the current value of the giving counter to
// giving boolean.
func (f *UintCounter) Add(b uint) {
	atomic.AddUint64(&f.fl, uint64(b))
}

// State returns the current state of the switch.
func (f *UintCounter) Read() uint {
	return uint(atomic.LoadUint64(&f.fl))
}
