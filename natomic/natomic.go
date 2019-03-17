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
func (f *BoolSwitch) State() bool {
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
	m  sync.Mutex
	fl int64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *IntSwitch) Flip(b int64) {
	f.m.Lock()
	atomic.StoreInt64(&b.fl, b)
	f.m.Unlock()
}

// State returns the current state of the switch.
func (f *IntSwitch) State() int64 {
	return atomic.LoadInt64(&b.fl)
}


// UintSwitch defines a concurrent safe bool switch.
//
// Once created, you must not copy a UintSwitch.
type UintSwitch struct {
	m  sync.Mutex
	fl uint64
}

// Flip flips the current value of the giving switch to
// giving boolean.
func (f *UintSwitch) Flip(b uint64) {
	f.m.Lock()
	atomic.StoreUint64(&b.fl, b)
	f.m.Unlock()
}

// State returns the current state of the switch.
func (f *UintSwitch) State() uint64 {
	return atomic.LoadUint64(&b.fl)
}