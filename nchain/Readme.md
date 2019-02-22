FutureChain
---------------
FutureChain implements a future-like inter-dependent pipeline using golang error group.
It allows usage of the power which the [Error Group](http://golang.org/x/sync/errgroup) package provides in safe conccurent 
functions with returned error.

## Install

```go
go get github.com/gokit/npkg/futurechain
```

## Example

- Sequential concurrent executions where one future is depent on the completion of 
another concurrent future operation.

```go
chain := futurechain.NewFutureChain(context.Background(), func() error {
	return nil
}).Go(func() error{
	// This will be executed with the chain function concurrently using 
	// goroutines.
	return nil
}).When(func() error {
	// First dependent sequentail concurrent chain.
	return nil
}).Go(func() error {
	// We be executed concurrently with the sequential concurrent chain 
	// returned by Chain.When().
	return nil
})

chain.Wait()
```

- Deferred future chains

This are cases where you can create a future which may be passed around to different functions
which will be executed by another future not yet created or ready.

```go
chain := futurechain.NewFutureChain(context.Background(), func() error {
	// ....
	return nil
})

// Deferred  chain can be created, if will not be executed until another 
// chain triggers it based on whatever method you use to chain.
chain2 := futurechain.DeferredChain(context.Background(), func() error {
	// ....
	return nil
})

// We want to trigger deffered chain regardless of error from first chain.
chain.ChainFuture(chain2)

chain2.Wait()
```
