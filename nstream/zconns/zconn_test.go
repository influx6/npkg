package zconns

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZConn(t *testing.T) {
	var ctx, cancel = context.WithCancel(context.Background())
	var conns = make(chan net.Conn)
	var errs = make(chan error, 1)

	var listener, err = ListenTCP(ctx, ":4050", errs, conns)
	require.NoError(t, err)
	require.NotNil(t, listener)

	defer cancel()

	var waiter sync.WaitGroup
	waiter.Add(1)

	go func() {
		defer waiter.Done()

		for conns := range conns {
			zc = NewZConn(conns)
			go func() {
				handleConnection(ctx, zc)
			}()
		}
	}()

	waiter.Wait()
}

func handleConnection(ctx context.Context, conn *ZConn) {
	reads := conn.Reads()
	writes := conn.Writes()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var content = bytes.NewBuffer(make([]byte, 512))
		payload := AcquireZPayload()
		payload.Stream = &nopWriter{content}

		reads <- payload
	}
}

type nopWriter struct {
	io.ReadWriter
}

func (nopWriter) Close() error {
	return nil
}
