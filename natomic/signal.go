package natomic

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gokit/npkg/nerror"
)

//************************************************************************
// Signals
//************************************************************************

// Signal defines a underline event representing some type and signal.
//
// The Type() return defines the unique means of identification
// for this signal.
type Signal interface {
	Type() string
}

// SignalResponder defines an interface for sending signals
// to it's implementer.
type SignalResponder interface {
	Respond(Signal)
}

// SignalResponderNotification defines an interface defining
// methods on a SignalResponder to deliver certain information
// regarding it's addition, removal from a ResponderGroup.
type SignalResponderNotification interface {
	SignalResponder

	OnAdded()
	OnRemoved()
}

// SignalReceiver defines a host of signal delivery and
// de-registration.
type SignalReceiver interface {
	SignalResponder

	// Add adds new SignalResponder to the underline signal receiver.
	// All signal provider to this SignalReceiver will be sent to the
	// SignalResponder as well.
	Add(SignalResponder) error

	// Remove must remove the responder from the underline
	// signal receiver.
	Remove(SignalResponder) error
}

//*****************************************************
// Atom
//*****************************************************

// Atom implements the Atom interface, implementing
// the safe concurrent storing and reading of stored values
// without the use of mutex and relying on the atomic.Value
// construct, which is great for low-write and high-read usage.
//
// Atoms are concurrently safe to use both for getting/setting
// value and adding interceptors.
type Atom struct {
	signals SignalReceiver
	imx     sync.Mutex
	store   *atomic.Value
}

// NewAtom returns a new instance of *Atom.
func NewAtom(receiver SignalReceiver) *Atom {
	var newStore atomic.Value
	return &Atom{
		signals: receiver,
		store:   &newStore,
	}
}

// Set attempts to set giving value into atom, if giving value is
// not the same underline type as previous set calls, then an error is
// returned.
func (am *Atom) Set(val interface{}) error {
	am.imx.Lock()
	defer am.imx.Unlock()
	am.store.Store(val)

	if am.signals != nil {
		am.signals.Respond(am)
	}
	return nil
}

// Type returns a giving string to represent the behaviour of
// giving type.
func (am *Atom) Type() string {
	return "natomic.Atom"
}

// Read returns the giving value stored within giving atom.
// It returns nil if no value was ever set.
func (am *Atom) Read() interface{} {
	return am.store.Load()
}

// Add adds a new SignalResponder into the atom responder
// provider.
func (am *Atom) Add(n SignalResponder) error {
	if am.signals != nil {
		return am.signals.Add(n)
	}
	return nil
}

// Remove removes giving SignalResponder from the atom responder
// provider.
func (am *Atom) Remove(n SignalResponder) error {
	if am.signals != nil {
		return am.signals.Remove(n)
	}
	return nil
}

// Respond delivers giving Signal to all interceptors
// registered with giving atom.
func (am *Atom) Respond(s Signal) {
	if am.signals != nil {
		am.signals.Respond(s)
	}
}

//************************************************************************
// ResponderGroup
//************************************************************************

// ResponderGroup defines a grouping and manager for SignalResponders, which manages
// the addition, and removals, basically the management of the lifecycle of a set of
// SignalResponders.
type ResponderGroup struct {
	sm         sync.Mutex
	spawned    bool
	waiter     sync.WaitGroup
	closer     chan struct{}
	actions    chan func()
	responders map[SignalResponder]struct{}
}

// NewResponderGroup returns a new ResponderGroup instance.
func NewResponderGroup() *ResponderGroup {
	var rg ResponderGroup
	rg.actions = make(chan func())
	rg.closer = make(chan struct{})
	rg.responders = map[SignalResponder]struct{}{}
	return &rg
}

// Remove returns an error if it fails to deliver SignalResponder into
// responder group.
func (rg *ResponderGroup) Remove(r SignalResponder) error {
	return rg.RemoveGuaranteed(r, nil)
}

// RemoveGuaranteed returns an error if it fails to deliver SignalResponder into
// responder group.
func (rg *ResponderGroup) RemoveGuaranteed(r SignalResponder, signal chan struct{}) error {
	var action = func() {
		if _, found := rg.responders[r]; found {
			delete(rg.responders, r)
			if irs, ok := r.(SignalResponderNotification); ok {
				irs.OnRemoved()
			}
		}

		if signal != nil {
			signal <- struct{}{}
		}
	}

	select {
	case rg.actions <- action:
		return nil
	case <-time.After(time.Second):
		return nerror.New("failed to remove responder")
	}
}

// Add returns an error if it fails to deliver SignalResponder into
// responder group.
func (rg *ResponderGroup) Add(r SignalResponder) error {
	return rg.AddGuaranteed(r, nil)
}

// AddGuaranteed returns an error if it fails to deliver SignalResponder into
// responder group.
func (rg *ResponderGroup) AddGuaranteed(r SignalResponder, signal chan struct{}) error {
	var action = func() {
		rg.responders[r] = struct{}{}
		if irs, ok := r.(SignalResponderNotification); ok {
			irs.OnAdded()
		}

		if signal != nil {
			signal <- struct{}{}
		}
	}

	select {
	case rg.actions <- action:
		return nil
	case <-time.After(time.Second):
		return nerror.New("failed to add responder")
	}
}

// Respond notifies giving signal to all responders.
func (rg *ResponderGroup) Respond(s Signal) {
	rg.RespondGuaranteed(s, nil)
}

// RespondGuaranteed returns an error if it fails to deliver SignalResponder into
// responder group.
func (rg *ResponderGroup) RespondGuaranteed(s Signal, signal chan struct{}) {
	var action = func() {
		for responder := range rg.responders {
			responder.Respond(s)
		}
		if signal != nil {
			signal <- struct{}{}
		}
	}

	rg.actions <- action
}

// Close closes giving responder group, removing all
// it's subscription.
func (rg *ResponderGroup) Close() {
	rg.closer <- struct{}{}
	rg.waiter.Wait()
}

// Start runs giving responder group.
//
// Start accepts possible many channels but takes the first one
// and will deliver a signal to it to provide detail that it
// as started.
func (rg *ResponderGroup) Start(started ...chan struct{}) {
	if rg.running() {
		return
	}

	if rg.actions == nil {
		rg.actions = make(chan func())
	}
	if rg.closer == nil {
		rg.closer = make(chan struct{})
	}
	if rg.responders == nil {
		rg.responders = map[SignalResponder]struct{}{}
	}

	if len(started) > 0 {
		rg.manage(started[0])
		return
	}

	rg.manage(nil)
}

// Wait blocks till giving ResponderGroup is closed.
func (rg *ResponderGroup) Wait() {
	rg.waiter.Wait()
}

func (rg *ResponderGroup) flip(t bool) {
	rg.sm.Lock()
	rg.spawned = t
	rg.sm.Unlock()
}

func (rg *ResponderGroup) running() bool {
	var running bool
	rg.sm.Lock()
	running = rg.spawned
	rg.sm.Unlock()
	return running
}

// manage setups the necessary goroutine for managing the underline
// process for giving responder.
//
// If a channel is passed then the giving channel receives a
// signal that giving goroutine as started.
func (rg *ResponderGroup) manage(started chan struct{}) {
	rg.flip(true)
	rg.waiter.Add(1)

	go func() {
		defer rg.waiter.Done()
		defer rg.flip(false)

		if started != nil {
			started <- struct{}{}
		}

		for {
			select {
			case <-rg.closer:
				rg.responders = map[SignalResponder]struct{}{}
				return
			case action := <-rg.actions:
				action()
			}
		}
	}()
}
