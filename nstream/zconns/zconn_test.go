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

func TestZConn(t *testing.T) {

	var listener, err = net.Listen("tcp", ":4050")
	require.NoError(t, err)
	require.NotNil(t, listener)

	var handler connHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener)
	server.Serve()

	var writeMessage = []byte("wondering through the ancient seas of the better world")

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

		<-payload.Done

		writePayload := AcquireZPayload()
		writePayload.Stream = payload.Stream

		writes <- writePayload

		<-writePayload.Done

		ReleaseZPayload(payload)
		ReleaseZPayload(writePayload)
		break
	}

	return nil
}

type nopWriter struct {
	io.ReadWriter
}

func (nopWriter) Close() error {
	return nil
}
