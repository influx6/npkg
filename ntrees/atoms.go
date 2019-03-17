package ntrees

import (
	"sync/atomic"

	"github.com/gokit/npkg/natomic"
)

// AtomSet defines the set method requirements for a safe
// concurrently usable implementer.
type AtomSet interface {
	Set(interface{}) error
}

// AtomRead defines the get method requirements for a safe
// concurrently usable implementer.
type AtomRead interface {
	Read() interface{}
}

// Atom exposes methods to safely set and get a giving
// underline value which can be safely retrieved atomically
// and concurrently.
type Atom interface {
	AtomSet
	AtomRead

	natomic.SignalResponder
}

// IndexUpdated defines a int type which is used
// to represent a index update.
type IndexUpdated int

// Type implements the Signal interface.
func (IndexUpdated) Type() string {
	return "INDEX_UPDATED"
}

// IntAtomImpl implements the Int interface, implementing
// the safe concurrent storing and reading of stored values
// without the use of mutex and relying on the atomic.Value
// construct, which is great for low-write and high-read usage.
type IntAtomImpl struct {
	v int64
}

// IntAtom returns a new instance of *IntAtomImpl.
func IntAtom() *IntAtomImpl {
	return &IntAtomImpl{v: -1}
}

// Set attempts to set giving value into atom, if giving value is
// not the same underline type as previous set calls, then an error is
// returned.
func (am *IntAtomImpl) Set(val int) {
	atomic.StoreInt64(&am.v, int64(val))
}

// Get returns the giving value stored within giving atom.
// It returns 0 if no value was ever set unless value set was 0.
func (am *IntAtomImpl) Read() int {
	return int(atomic.LoadInt64(&am.v))
}
