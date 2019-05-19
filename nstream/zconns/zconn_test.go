package zconns

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	message = []byte("wondering through the ancient seas of the better world endless awaiting the pranch")
)

func BenchmarkZConn(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.StopTimer()
	var listener, err = net.Listen("tcp", ":5050")
	if err != nil {
		panic(err)
	}

	var handler readHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener, false)
	server.Serve()

	var clientConn, clientErr = net.DialTimeout("tcp", ":5050", time.Second*5)
	if clientErr != nil {
		panic(clientErr)
	}

	var readBufferSize = 4096
	var zclient = NewZConn(clientConn, ZConnParentContext(ctx), ZConnReadBuffer(readBufferSize), ZConnWriteBuffer(readBufferSize))

	var readBuffer = bytes.NewReader(message)
	var readContent = ioutil.NopCloser(readBuffer)

	go func() {
		_, _ = io.Copy(ioutil.Discard, clientConn)
	}()

	b.SetBytes(int64(len(message)))
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		readBuffer.Reset(message)
		if err := zclient.WriteTo(readContent, false); err != nil {
			panic(err)
		}
	}

	_ = zclient.Flush()

	b.StopTimer()
	cancel()

	_ = zclient.Close()
	_ = server.Wait()
}

func TestZConn(t *testing.T) {
	var listener, err = net.Listen("tcp", ":4050")
	require.NoError(t, err)
	require.NotNil(t, listener)

	var handler readHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener, true)
	server.Serve()

	var clientConn, clientErr = net.DialTimeout("tcp", ":4050", time.Second*5)
	require.NoError(t, clientErr)
	require.NotNil(t, clientConn)

	var zclient = NewZConn(clientConn, ZConnDebugMode())

	var writeContent = noCloser(bytes.NewBuffer(message))
	require.NoError(t, zclient.WriteTo(writeContent, true))

	var readBuffer = bytes.NewBuffer(make([]byte, 0, 512))
	var readContent = noCloser(readBuffer)
	require.NoError(t, zclient.ReadFrom(readContent))
	require.Equal(t, message, readBuffer.Bytes())
	require.NoError(t, zclient.Close())

	cancel()
	require.Error(t, server.Wait())
}

type readHandler struct{}

func (readHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	var zc = NewZConn(conn, ZConnParentContext(ctx))
	var buffer = bytes.NewBuffer(make([]byte, 0, 512))
	var writeContent = noCloser(buffer)

	for {
		select {
		case <-ctx.Done():
			_ = zc.Close()
			return nil
		default:
		}

		buffer.Reset()

		if err := zc.ReadFrom(writeContent); err != nil {
			log.Printf("[ConnHandler] | %s | Closing serverConn due to read error", zc.id)
			return err
		}
	}
}

type writeHandler struct{}

func (writeHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	var zc = NewZConn(conn, ZConnParentContext(ctx))

	var writeBuffer = bytes.NewReader(message)
	var writeContent = ioutil.NopCloser(writeBuffer)

	for {
		select {
		case <-ctx.Done():
			_ = zc.Close()
			return nil
		default:
		}

		writeBuffer.Reset(message)

		if err := zc.WriteTo(writeContent, true); err != nil {
			log.Printf("[ConnHandler] | %s | Closing serverConn due to write error", zc.id)
			return err
		}
	}
}

func noCloser(w io.ReadWriter) io.ReadWriteCloser {
	return &nopCloser{w}
}

type nopCloser struct {
	io.ReadWriter
}

func (nopCloser) Close() error {
	return nil
}
