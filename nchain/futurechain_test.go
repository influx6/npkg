package nchain_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gokit/npkg/nchain"
	"github.com/stretchr/testify/assert"
)

func ExampleFutureChain_When() {
	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return nil
	}).When(func(_ context.Context) error {
		return nil
	})

	chain.Wait()
}

func ExampleFutureChain_Then() {
	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return nil
	}).Then(func(_ context.Context) error {
		return nil
	})

	chain.Wait()
}

func ExampleFutureChain_Chain() {
	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return nil
	}).Chain(func(_ context.Context) error {
		return nil
	})

	chain.Wait()
}

func TestNewFutureChain(t *testing.T) {
	signal := make(chan struct{}, 2)
	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		signal <- struct{}{}
		return nil
	})

	assert.NoError(t, chain.Wait())
	assert.Len(t, signal, 1)
}

func TestFutureChain_Call(t *testing.T) {
	var firstChainTime time.Time
	var secondChainTime time.Time

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		firstChainTime = time.Now()
		return nil
	}).Chain(func(_ context.Context) error {
		secondChainTime = time.Now()
		return nil
	})

	assert.NoError(t, chain.Wait())
	assert.True(t, firstChainTime.Before(secondChainTime))
}

func TestFutureChain_FutureCall(t *testing.T) {
	var firstChainTime time.Time
	var secondChainTime time.Time

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		firstChainTime = time.Now()
		return nil
	})

	chain2 := nchain.DeferredChain(context.Background(), func(_ context.Context) error {
		secondChainTime = time.Now()
		return nil
	})

	chain.ChainFuture(chain2)

	assert.NoError(t, chain2.Wait())
	assert.True(t, firstChainTime.Before(secondChainTime))
}

func TestFutureChain_FutureChain(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	})

	chain2 := nchain.DeferredChain(context.Background(), func(_ context.Context) error {
		return secondError
	})

	chain.ChainFuture(chain2)

	assert.Error(t, chain2.Wait())
	assert.Equal(t, chain2.Wait(), secondError)
}

func TestFutureChain_Chain(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	}).Chain(func(_ context.Context) error {
		return secondError
	})

	assert.Error(t, chain.Wait())
	assert.Equal(t, chain.Wait(), secondError)
}

func TestFutureChain_Then(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	})

	chain2 := chain.Chain(func(_ context.Context) error {
		return secondError
	})

	assert.Error(t, chain.Wait())
	assert.Equal(t, chain.Wait(), firstError)
	assert.Equal(t, chain2.Wait(), firstError)
}

func TestFutureChain_ThenFuture_MultiChain(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	})

	chain2 := nchain.DeferredChain(context.Background(), func(_ context.Context) error {
		return secondError
	})

	chain2.Go(func(_ context.Context) error {
		return nil
	})

	chain2.Go(func(_ context.Context) error {
		return nil
	})

	chain.ThenFuture(chain2)

	assert.Error(t, chain2.Wait())
	assert.Equal(t, chain2.Wait(), firstError)
}

func TestFutureChain_ThenFuture(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	})

	chain2 := nchain.DeferredChain(context.Background(), func(_ context.Context) error {
		return secondError
	})

	chain.ThenFuture(chain2)

	assert.Error(t, chain2.Wait())
	assert.Equal(t, chain2.Wait(), firstError)
}

func TestFutureChain_WhenFuture(t *testing.T) {
	firstError := errors.New("first error")
	secondError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return secondError
	})

	signal := make(chan struct{}, 1)
	chain2 := nchain.DeferredChain(context.Background(), func(_ context.Context) error {
		signal <- struct{}{}
		return firstError
	})

	chain.WhenFuture(chain2)

	assert.Error(t, chain2.Wait())
	assert.Equal(t, chain2.Wait(), secondError)
	assert.Len(t, signal, 0)
}

func TestFutureChain_When(t *testing.T) {
	firstError := errors.New("first error")

	chain := nchain.NewFutureChain(context.Background(), func(_ context.Context) error {
		return firstError
	})

	signal := make(chan struct{}, 1)
	secondChain := chain.When(func(_ context.Context) error {
		signal <- struct{}{}
		return nil
	})

	assert.Error(t, secondChain.Wait())
	assert.Equal(t, secondChain.Wait(), firstError)
	assert.Len(t, signal, 0)
}
