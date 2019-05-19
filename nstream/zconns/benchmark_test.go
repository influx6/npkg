package zconns

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"testing"
	"time"
)

func BenchmarkZConnMessages(b *testing.B) {
	b.Run("0 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(0)
		benchZConn(b, payload)
	})

	b.Run("2 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(2)
		benchZConn(b, payload)
	})

	b.Run("4 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(4)
		benchZConn(b, payload)
	})

	b.Run("8 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(8)
		benchZConn(b, payload)
	})

	b.Run("16 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(16)
		benchZConn(b, payload)
	})

	b.Run("32 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(32)
		benchZConn(b, payload)
	})

	b.Run("64 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(64)
		benchZConn(b, payload)
	})

	b.Run("128 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(128)
		benchZConn(b, payload)
	})

	b.Run("256 Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(256)
		benchZConn(b, payload)
	})

	b.Run("1K Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(1024)
		benchZConn(b, payload)
	})

	b.Run("2K Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(1024 * 2)
		benchZConn(b, payload)
	})

	b.Run("4K Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(1024 * 4)
		benchZConn(b, payload)
	})

	b.Run("8K Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(1024 * 8)
		benchZConn(b, payload)
	})

	b.Run("16K Bytes", func(b *testing.B) {
		b.StopTimer()
		payload := sizedPayload(1024 * 16)
		benchZConn(b, payload)
	})
}

func benchZConn(b *testing.B, message []byte) {
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
		if err := zclient.ReadFrom(readContent, false); err != nil {
			panic(err)
		}
	}

	_ = zclient.Flush()

	b.StopTimer()
	cancel()

	_ = zclient.Close()
	_ = server.Wait()
}

var pub = "PUB "
var ch = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@$#%^&*()")

func sizedPayloadString(sz int) string {
	return string(sizedPayload(sz))
}

func sizedPayload(sz int) []byte {
	payload := make([]byte, len(pub)+sz)
	nx := copy(payload, pub)
	copy(payload[nx:], sizedBytes(sz))
	return payload
}

func sizedBytes(sz int) []byte {
	if sz <= 0 {
		return []byte("")
	}

	b := make([]byte, sz)
	for i := range b {
		b[i] = ch[rand.Intn(len(ch))]
	}
	return b
}

func sizedString(sz int) string {
	return string(sizedBytes(sz))
}
