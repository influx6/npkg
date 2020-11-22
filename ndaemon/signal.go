package ndaemon

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func WaitForKill() os.Signal {
	var c = WaitForKillChan()
	return <-c
}

func WaitForKillChan() chan os.Signal {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, os.Kill)
	return interrupt
}

// WaiterForKillWithSignal will call the canceler function when a interrupt/kill signal is received.
func WaiterForKillWithSignal(signalChan chan os.Signal, canceler context.CancelFunc) *sync.WaitGroup {
	var waiter sync.WaitGroup
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		<-signalChan
		canceler()
	}()
	return &waiter
}

// WaiterForCtxSignal returns a new sync.WaitGroup which will be completed
// once either the context gets closed due to another goroutine closing the
// cancel function or the signal from the os to kill the operation.
func WaiterForCtxSignal(ctx context.Context, canceler context.CancelFunc) *sync.WaitGroup {
	var signalChan = WaitForKillChan()
	var waiter sync.WaitGroup
	waiter.Add(1)
	go func() {
		defer waiter.Done()
		defer canceler()
		select {
		case <-ctx.Done():
			break
		case <-signalChan:
			break
		}
	}()
	return &waiter
}

// CtxAndWaiterForSignal returns a new context.Context and WaitGroup which can be
// depended on to be closed once the os kill/interrupt signal is received.
func CtxAndWaiterForSignal(ctx context.Context, signalChan chan os.Signal) (context.Context, *sync.WaitGroup) {
	var newCtx, canceler = context.WithCancel(ctx)
	var waiter = WaiterForKillWithSignal(signalChan, canceler)
	return newCtx, waiter
}

// CtxAndWaiterFor returns a new context.Context and WaitGroup which can be
// depended on to be closed once the os kill/interrupt signal is received.
func CtxAndWaiterFor(ctx context.Context) (context.Context, *sync.WaitGroup) {
	var newCtx, canceler = context.WithCancel(ctx)
	var signalChan = WaitForKillChan()
	var waiter = WaiterForKillWithSignal(signalChan, canceler)
	return newCtx, waiter
}
