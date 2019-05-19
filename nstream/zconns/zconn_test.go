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
	require.NoError(t, zclient.ReadFrom(writeContent, true))

	var readBuffer = bytes.NewBuffer(make([]byte, 0, 512))
	var readContent = noCloser(readBuffer)
	require.NoError(t, zclient.WriteTo(readContent))
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

		if err := zc.WriteTo(writeContent); err != nil {
			//log.Printf("[readHandler] | %s | Closing serverConn due to read error", zc.id)
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

		if err := zc.ReadFrom(writeContent, true); err != nil {
			log.Printf("[writeHandler] | %s | Closing serverConn due to write error", zc.id)
			return err
		}
	}
}

type readWriteHandler struct{}

func (readWriteHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	var zc = NewZConn(conn, ZConnParentContext(ctx))
	var buffer = bytes.NewBuffer(make([]byte, 0, 512))
	var writeContent = noCloser(buffer)
	var readCloser = ioutil.NopCloser(buffer)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		buffer.Reset()

		if err := zc.WriteTo(writeContent); err != nil {
			log.Printf("[readWriteHandler] | %s | Closing serverConn due to read error", zc.id)
			return err
		}

		if err := zc.ReadFrom(readCloser, true); err != nil {
			log.Printf("[readWriteHandler] | %s | Closing serverConn due to write error", zc.id)
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
