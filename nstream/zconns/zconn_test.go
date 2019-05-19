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

	var handler connHandler
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
		if err := zclient.Write(readContent, false); err != nil {
			panic(err)
		}
	}

	_ = zclient.Flush()

	b.StopTimer()
	cancel()

	log.Println("Closing....")
	_ = zclient.Close()
	log.Println("Waiting....")
	_ = server.Wait()
	log.Println("Finished Waiting....")
}

func TestZConn(t *testing.T) {
	var listener, err = net.Listen("tcp", ":4050")
	require.NoError(t, err)
	require.NotNil(t, listener)

	var handler connHandler
	var ctx, cancel = context.WithCancel(context.Background())
	var server = NewServer(ctx, handler, listener, true)
	server.Serve()

	var clientConn, clientErr = net.DialTimeout("tcp", ":4050", time.Second*5)
	require.NoError(t, clientErr)
	require.NotNil(t, clientConn)

	var zclient = NewZConn(clientConn, ZConnDebugMode())

	var writeContent = noCloser(bytes.NewBuffer(message))
	require.NoError(t, zclient.Write(writeContent, true))

	var readBuffer = bytes.NewBuffer(make([]byte, 0, 512))
	var readContent = noCloser(readBuffer)
	require.NoError(t, zclient.Read(readContent, true))
	require.Equal(t, message, readBuffer.Bytes())
	require.NoError(t, zclient.Close())

	cancel()
	require.Error(t, server.Wait())
}

type connHandler struct{}

func (connHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	defer log.Printf("[ConnHandler] | Closing serverConn")

	var zc = NewZConn(conn, ZConnParentContext(ctx))
	var buffer = bytes.NewBuffer(make([]byte, 0, 512))
	var writeContent = noCloser(buffer)

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing connHandler client....")
			_ = zc.Close()
			log.Println("Closing connHandler....")
			return nil
		default:
		}

		buffer.Reset()

		if err := zc.Read(writeContent, true); err != nil {
			log.Printf("[ConnHandler] | %s | Closing serverConn due to read error", zc.id)
			return err
		}

		//if err := zc.Write(writeContent, true); err != nil {
		//	log.Printf("[ConnHandler] | %s | Closing serverConn due to write error", zc.id)
		//	return err
		//}
	}
}

type writeConnHandler struct{}

func (writeConnHandler) ServeConn(ctx context.Context, conn net.Conn) error {
	var zc = NewZConn(conn, ZConnParentContext(ctx))

	var writeBuffer = bytes.NewBuffer(message)
	var writeContent = noCloser(writeBuffer)
	var readBuffer = bytes.NewBuffer(make([]byte, 0, 512))
	var readContent = noCloser(readBuffer)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[ConnHandler] | %s | Closing serverConn", zc.id)
			_ = zc.Close()
			return nil
		default:
		}

		readBuffer.Reset()

		if err := zc.Write(writeContent, true); err != nil {
			log.Printf("[ConnHandler] | %s | Closing serverConn due to write error", zc.id)
			return err
		}

		if err := zc.Read(readContent, true); err != nil {
			log.Printf("[ConnHandler] | %s | Closing serverConn due to read error", zc.id)
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
