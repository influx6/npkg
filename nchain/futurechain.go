package nchain

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// FutureChain implements a pure Future implementation which
// will lunch an operation in a go-routine but block till
// that operation finishes. It is created to allow
// chaining next operation over a linear call stack.
//
// FutureChain relies on errgroup.Group underneath.
//
type FutureChain struct {
	ctx context.Context
	g   *errgroup.Group
	fn  func(context.Context) error

	error     chan error
	closer    chan struct{}
	signal    chan struct{}
	revSignal chan struct{}
}

func noWork(_ context.Context) error { return nil }

// NoWorkChain returns a new FutureChain which has returns nil immediately
// for it's work function. This allows creating new chains that already are
// resolved to be used to chain other incoming work.
func NoWorkChain(ctx context.Context) *FutureChain {
	return NewFutureChain(ctx, noWork)
}

// NewFutureChain returns a new instance of a FutureChain using provided function.
// It immediately lunches function into future chain error group.
func NewFutureChain(ctx context.Context, fx func(context.Context) error) *FutureChain {
	chain := newFutureChain(ctx, fx)
	chain.exec(nil, true)
	return chain
}

// DeferredChain returns an new Future where it's operation will not be started
// immediately but will be initiated by another completed Future chain.
// This is useful if you wish to reduce callback trees and instead have
// specific instantiated FutureChains that can be passed to one another.
// This is also useful for deferred calculations that will be performed
// based on the completion of some other future yet to be created.
func DeferredChain(ctx context.Context, fx func(context.Context) error) *FutureChain {
	chain := newFutureChain(ctx, fx)
	return chain
}

// NewFutureChain returns a new instance of a FutureChain using provided function.
func newFutureChain(ctx context.Context, fx func(context.Context) error) *FutureChain {
	ew, ctx := errgroup.WithContext(ctx)
	chain := &FutureChain{
		g:         ew,
		ctx:       ctx,
		fn:        fx,
		error:     make(chan error, 0),
		closer:    make(chan struct{}, 0),
		signal:    make(chan struct{}, 0),
		revSignal: make(chan struct{}, 0),
	}

	// schedule error kicker.
	chain.g.Go(func() error {
		defer func() {
			chain.revSignal <- struct{}{}
		}()
		return <-chain.error
	})

	// schedule function execution.
	chain.Go(fx)

	return chain
}

// Go just adds another function into this current futures
// wait group, it behaves exactly like the errgroup.Group.Go
// method.
//
// It returns itself, to allow chaining.
func (f *FutureChain) Go(fx func(context.Context) error) *FutureChain {
	f.g.Go(func() error {
		<-f.closer
		select {
		case <-f.signal:
			return fx(f.ctx)
		default:
			return nil
		}
	})
	return f
}

// When chains the next function call when this future has completed
// without having an error, else will just pass occurred error down
// the chain.
//
// It always returns a new FutureChain instance.
//
func (f *FutureChain) When(fx func(context.Context) error) *FutureChain {
	newChain := newFutureChain(f.ctx, fx)
	newChain.execIfNoError(f)
	return newChain
}

// WhenFuture maps giving FutureChain to be called when this future
// has executed successfully with no error else, passing down occurred error.
//
// It always returns its self for chaining.
//
func (f *FutureChain) WhenFuture(fm *FutureChain) *FutureChain {
	f.execChildIfNoError(fm)
	return f
}

// Then chains the next function call regardless of the failure of
// this future in the chain. It will get executed and the giving
// occurred error from this chain will be passed down the chain
// regardless of any error occurring from the function passed to the
// returned future (i.e it always returns parent's error down the chain).
//
// It always returns a new FutureChain instance.
//
func (f *FutureChain) Then(fx func(context.Context) error) *FutureChain {
	newChain := newFutureChain(f.ctx, fx)
	newChain.execBranchDownError(f)
	return newChain
}

// ThenFuture maps giving FutureChain to be called when this future
// has executed successfully with or without an error, if there is an
// error then the error is passed down to the provided future but the
// provided future will still be executed.
//
// It always returns its self for chaining.
//
func (f *FutureChain) ThenFuture(fm *FutureChain) *FutureChain {
	f.execChildEvenIfErrorWithPassDown(fm)
	return f
}

// Chain chains the next function call regardless of the failure of
// this future in the chain. It will get executed after completion of
// this future and the giving occurred error from this chain will be ignored.
// This forces you to keep a reference to this future chain to retrieve the
// error that occurred for this chain as the returned future will only ever
// return an error that occurred from calling the function passed  in.
func (f *FutureChain) Chain(fx func(context.Context) error) *FutureChain {
	newChain := newFutureChain(f.ctx, fx)
	newChain.execNoBranchDownError(f)
	return newChain
}

// ChainFuture maps giving FutureChain to be called when this future
// has executed successfully with or without an error, if there is an
// error then the error is still not passed down to the provided future
// and but the provided future will be executed.
//
// It always returns its self for chaining.
//
func (f *FutureChain) ChainFuture(fm *FutureChain) *FutureChain {
	f.execChildEvenIfError(fm)
	return f
}

// Wait blocks till the operation is completed and
// returns error seen.
func (f *FutureChain) Wait() error {
	return f.g.Wait()
}

// exec launches function within future chain error group.
func (f *FutureChain) exec(err error, execFunc bool) {
	f.error <- err
	<-f.revSignal
	if execFunc {
		close(f.signal)
	}
	close(f.closer)
}

// execIfNoError will execute future chains function if parent chain
// has no error occurred else returns parent error down the chain.
//
// It will wait until parent future has finished execution.
func (f *FutureChain) execIfNoError(parent *FutureChain) {
	go func() {
		err := parent.Wait()
		f.exec(err, err == nil)
	}()
}

// execChildEvenIfError will execute provided chain if giving chain
// has no error.
func (f *FutureChain) execChildIfNoError(child *FutureChain) {
	go func() {
		err := f.Wait()
		child.exec(err, err == nil)
	}()
}

// execChildEvenIfError will execute provided chain even if giving chain
// had an error but will not pass down parents error.
func (f *FutureChain) execChildEvenIfError(child *FutureChain) {
	go func() {
		f.Wait()
		child.exec(nil, true)
	}()
}

// execChildEvenIfErrorWithPassDown will execute provided chain even if giving chain
// had an error and will pass down parents error.
func (f *FutureChain) execChildEvenIfErrorWithPassDown(child *FutureChain) {
	go func() {
		parentErr := f.Wait()
		child.exec(parentErr, true)
	}()
}

// execBranchDownError will execute future chains function even if parent chain
// has an error but it returns parent error down the chain even if chain function
// had an error occur.
//
// It will wait until parent future has finished execution.
func (f *FutureChain) execBranchDownError(parent *FutureChain) {
	go func() {
		err := parent.Wait()
		f.exec(err, true)
	}()
}

// execNoBranchDownError will execute future chains function even if parent chain
// has an error, it will only ever pass down error from child function.
//
// It will wait until parent future has finished execution.
func (f *FutureChain) execNoBranchDownError(parent *FutureChain) {
	go func() {
		parent.Wait()
		f.exec(nil, true)
	}()
}
