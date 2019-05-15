package zconns

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	writeMessage = []byte("wondering through the ancient seas of the better world")
)

func BenchmarkZConn(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var listener, err = net.Listen("tcp", ":5050")
	if err != nil {
		panic(err)
	}

	var handler connHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener)
	server.Serve()

	b.SetBytes(int64(len(writeMessage)))

	var clientConn, clientErr = net.DialTimeout("tcp", ":5050", time.Second*5)
	if clientErr != nil {
		panic(clientErr)
	}

	var zclient = NewZConn(clientConn)
	writes := zclient.Writes()
	reads := zclient.Reads()

	var readPayload = AcquireZPayload()
	var readContent = bytes.NewBuffer(make([]byte, 0, 512))
	readPayload.Stream = &nopWriter{readContent}

	zclient.Reads() <- readPayload

	for i := 0; i < b.N; i++ {
		readContent.Reset()

		var buffer = bytes.NewBuffer(writeMessage)
		var nopbuffer = &nopWriter{buffer}
		var payload = AcquireZPayload()
		payload.Stream = nopbuffer

		writes <- payload
		<-payload.Done
		ReleaseZPayload(payload)

		reads <- readPayload
		<-readPayload.Done
	}

	zclient.Close()
	b.StopTimer()

	cancel()
	_ = server.Wait()
}

func TestZConn(t *testing.T) {
	var listener, err = net.Listen("tcp", ":4050")
	require.NoError(t, err)
	require.NotNil(t, listener)

	var handler connHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener)
	server.Serve()

	var clientConn, clientErr = net.DialTimeout("tcp", ":4050", time.Second*5)
	require.NoError(t, clientErr)
	require.NotNil(t, clientConn)

	var zclient = NewZConn(clientConn)

	payload := AcquireZPayload()
	var content = bytes.NewBuffer(writeMessage)
	payload.Stream = &nopWriter{content}

	zclient.Writes() <- payload

	select {
	case err := <-payload.Err:
		require.Fail(t, "Failed with write: %s", err)
	case <-payload.Done:
	}

	var readPayload = AcquireZPayload()
	var readContent = bytes.NewBuffer(make([]byte, 0, 512))
	readPayload.Stream = &nopWriter{readContent}

	zclient.Reads() <- readPayload

	select {
	case err := <-readPayload.Err:
		require.Fail(t, "Failed with read: %s", err)
	case <-readPayload.Done:
	}

	require.Equal(t, writeMessage, readContent.Bytes())

	require.NoError(t, zclient.Close())

	cancel()
	require.Error(t, server.Wait())
}

type connHandler struct{}

func (connHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	var zc = NewZConn(conn, ZConnParentContext(ctx))
	reads := zc.Reads()
	writes := zc.Writes()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var content = bytes.NewBuffer(make([]byte, 0, 512))
		payload := AcquireZPayload()
		payload.Stream = &nopWriter{content}

		reads <- payload

		select {
		case <-payload.Done:
		case err := <-payload.Err:
			ReleaseZPayload(payload)
			return err
		}

		writePayload := AcquireZPayload()
		writePayload.Stream = payload.Stream

		writes <- writePayload

		select {
		case <-writePayload.Done:
			continue
		case err := <-writePayload.Err:
			ReleaseZPayload(payload)
			ReleaseZPayload(writePayload)
			return err
		}

		ReleaseZPayload(payload)
		ReleaseZPayload(writePayload)
	}
}

type nopWriter struct {
	io.ReadWriter
}

func (nopWriter) Close() error {
	return nil
}
