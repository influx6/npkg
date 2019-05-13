package zconns

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gokit/npkg/nbytes"
)

//*********************************************************************************************
// ZTimedConn
//*********************************************************************************************

var (
	noTime = time.Time{}
)

// ZTimeout defines a function which takes previous timeout returning same or
// new timeout to use for next operation.
type ZTimeout func(time.Duration) time.Duration

// ZTimedConn implements a wrapper around a net.Conn which guards giving connection
// with appropriate read/write timeout.
type ZTimedConn struct {
	net.Conn
	readTimeout  ZTimeout
	writeTimeout ZTimeout
	defaultRead  time.Duration
	defaultWrite time.Duration
	lastRead     time.Duration
	lastWrite    time.Duration
}

// NewZTimedConn returns a new instance of a ZTimedConn.
func NewZTimedConn(conn net.Conn, defaultRead time.Duration, defaultWrite time.Duration, rd ZTimeout, wd ZTimeout) *ZTimedConn {
	return &ZTimedConn{
		Conn:         conn,
		readTimeout:  rd,
		writeTimeout: wd,
		defaultRead:  defaultRead,
		defaultWrite: defaultWrite,
	}
}

// Write calls the underline connection read with provided timeout.
func (c *ZTimedConn) Write(b []byte) (int, error) {
	c.lastWrite = c.readTimeout(c.lastWrite)
	var writeErr = c.Conn.SetWriteDeadline(time.Now().Add(c.lastWrite))
	if writeErr != nil {
		return 0, writeErr
	}

	var writeCount, err = c.Conn.Write(b)
	if err != nil {
		_ = c.Conn.SetWriteDeadline(noTime)
		return writeCount, err
	}

	var resetErr = c.Conn.SetWriteDeadline(noTime)
	if resetErr != nil {
		return writeCount, resetErr
	}

	c.lastWrite = c.defaultWrite
	return writeCount, nil
}

// Read calls the underline connection read with provided timeout.
func (c *ZTimedConn) Read(b []byte) (int, error) {
	c.lastRead = c.readTimeout(c.lastRead)
	var readErr = c.Conn.SetReadDeadline(time.Now().Add(c.lastRead))
	if readErr != nil {
		return 0, readErr
	}

	var readCount, err = c.Conn.Read(b)
	if err != nil {
		_ = c.Conn.SetReadDeadline(noTime)
		return readCount, err
	}

	var resetErr = c.Conn.SetReadDeadline(noTime)
	if resetErr != nil {
		return readCount, resetErr
	}

	c.lastRead = c.defaultRead
	return readCount, nil
}

//**********************************************************************
// ListenTCP
//**********************************************************************

// ListenTCP starts a tcp server listening on the specified address (default 127.0.0.1:1445)
// and then continually reads from the server handling any incoming connections
func ListenTCP(ctx context.Context, addr string, errChan chan<- error, connChan chan net.Conn) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = ln.Close()
		close(connChan)
	}()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				errChan <- fmt.Errorf("failed to accept TCP connection %s", err.Error())
				return
			}

			connChan <- conn
		}
	}()

	return ln, nil
}

//*********************************************************************************************
// Next Message Format: NX {MAX_64_INT}\r\n\r\n
// Data Message Format: DX {HEADER} {DATA_BYTES}\r\n\r\n
// Delimiter Format: \r\n
//**********************************************************************************************

var (
	// ErrKillConnection defines a returnable error from a
	// ZConnReader and ZConnWriter to ensure giving connection
	// get's cosed due to some underline error.
	ErrKillConnection = errors.New("unrecoverable error")
)

// ZApply defines a function type to apply a change to the ZConn instance.
type ZApply func(*ZConn)

// ZConnDefaultWriteTimeout sets the default timeout duration to be used for write
// calls for a ZConn.
func ZConnDefaultWriteTimeout(d time.Duration) ZApply {
	return func(conn *ZConn) {
		conn.defaultWriteTimeout = d
	}
}

// ZConnDefaultReadTimeout sets the default timeout duration to be used for read
// calls for a ZConn.
func ZConnDefaultReadTimeout(d time.Duration) ZApply {
	return func(conn *ZConn) {
		conn.defaultReadTimeout = d
	}
}

// ZConnWriteTimeout sets the ZTimeout function type to be used for write
// calls for a ZConn.
func ZConnWriteTimeout(fn ZTimeout) ZApply {
	return func(conn *ZConn) {
		conn.writeTimeouts = fn
	}
}

// ZConnReadTimeout sets the ZTimeout function type to be used for read
// calls for a ZConn.
func ZConnReadTimeout(fn ZTimeout) ZApply {
	return func(conn *ZConn) {
		conn.readTimeouts = fn
	}
}

// ZConnWriteBuffer sets the buffer space to be used for write
// calls for a ZConn.
func ZConnWriteBuffer(buffer int64) ZApply {
	return func(conn *ZConn) {
		conn.writeBuffer = buffer
	}
}

// ZConnReadBuffer sets the buffer space to be used for read
// calls for a ZConn.
func ZConnReadBuffer(buffer int64) ZApply {
	return func(conn *ZConn) {
		conn.readBuffer = buffer
	}
}

// ZConnMaxWrite sets the maximum allowed read size for each data stream.
// The maximum allowed size for each unique message.
func ZConnMaxRead(max int64) ZApply {
	return func(conn *ZConn) {
		conn.maxRead = max
	}
}

// ZConnMaxWrite sets the maximum allowed data size for collected data stream.
// The maximum allowed size for each unique message.
func ZConnMaxWrite(max int64) ZApply {
	return func(conn *ZConn) {
		conn.maxWrite = max
	}
}

// ZConnWriter defines the underline function to be used for
// handling write operation for a ZPayload.
type ZConnWriter func(io.Writer, *ZPayload) error

// ZConnReader defines the underline function to be used for
// handling read operation for a ZPayload.
type ZConnReader func(io.Reader, *ZPayload) error

// ZConnWriteWorker sets the worker function to be used for write
// calls for a ZConn.
func ZConnWriteWorker(worker ZConnWriter) ZApply {
	return func(conn *ZConn) {
		conn.writeWorker = worker
	}
}

// ZConnReadWorker sets the worker function to be used for read
// calls for a ZConn.
func ZConnReadWorker(worker ZConnReader) ZApply {
	return func(conn *ZConn) {
		conn.readWorker = worker
	}
}

// ZConnParentContext sets the parent context for giving
// ZConn to respect cancellation.
func ZConnParentContext(ctx context.Context) ZApply {
	return func(conn *ZConn) {
		conn.ctx, conn.ctxCanceler = context.WithCancel(ctx)
	}
}

// ZConnReadRequests sets a stream object to be used for read
// requests by a ZConn.
func ZConnReadRequests(stream ZPayloadStream) ZApply {
	return func(conn *ZConn) {
		conn.readRequests = stream
	}
}

// ZConnWriteRequests sets a stream object to be used for write
// requests by a ZConn.
func ZConnWriteRequests(stream ZPayloadStream) ZApply {
	return func(conn *ZConn) {
		conn.writeRequests = stream
	}
}

var (
	zPayloads = sync.Pool{
		New: func() interface{} {
			var payload = new(ZPayload)
			payload.Err = make(chan error, 1)
			payload.Done = make(chan struct{}, 1)
			return payload
		},
	}
)

// AcquireZPayload returns a new ZPayload for use.
func AcquireZPayload() *ZPayload {
	if nzp, ok := zPayloads.Get().(*ZPayload); ok {
		return nzp
	}
	return new(ZPayload)
}

// ReleaseZPayload returns giving ZPayload into underline object pool
// for reuse.
// Never re-use a ZPayload reference you've previously released.
func ReleaseZPayload(zp *ZPayload) {
	zPayloads.Put(zp)
}

// ZPayload defines an underline structure for writing data
// into an underline ZConn.
type ZPayload struct {
	Done   chan struct{}
	Err    chan error
	Addr   chan net.Addr
	Stream io.ReadWriteCloser
}

// Reset giving ZPayload fields.
func (z *ZPayload) Reset() {
	z.Err = nil
	z.Addr = nil
	z.Stream = nil
	z.Done = nil
}

// ZPayloadStream defines an underline channel type which represent
// a channel for the writing and reading of ZPayload.
type ZPayloadStream chan *ZPayload

// ZConn implements a simple connection wrapper around a
// net.Conn.
type ZConn struct {
	laddr               net.Addr
	addr                net.Addr
	maxRead             int64
	maxWrite            int64
	readBuffer          int64
	writeBuffer         int64
	readTimeouts        ZTimeout
	writeTimeouts       ZTimeout
	conn                net.Conn
	readWorker          ZConnReader
	writeWorker         ZConnWriter
	readRequests        ZPayloadStream
	writeRequests       ZPayloadStream
	waiter              sync.WaitGroup
	ctx                 context.Context
	ctxCanceler         context.CancelFunc
	defaultWriteTimeout time.Duration
	defaultReadTimeout  time.Duration
	streamWriter        *nbytes.DelimitedStreamWriter
	streamReader        *nbytes.DelimitedStreamReader
}

const (
	defaultTimeout     = time.Second * 5
	defaultReadBuffer  = 1024 * 4
	defaultWriteBuffer = 1024 * 4
	defaultEscape      = "/&"
	defaultDelimiter   = "\r\n"
)

// NewZConn returns a new instance of a ZConn.
func NewZConn(conn net.Conn, fns ...ZApply) *ZConn {
	var zc = new(ZConn)
	zc.conn = conn
	zc.addr = conn.RemoteAddr()
	zc.laddr = conn.LocalAddr()
	zc.readBuffer = defaultReadBuffer
	zc.writeBuffer = defaultWriteBuffer
	zc.defaultWriteTimeout = defaultTimeout
	zc.defaultWriteTimeout = defaultTimeout

	var defaultTimer = sameTimeout(defaultTimeout)
	zc.readTimeouts = defaultTimer
	zc.readTimeouts = defaultTimer

	for _, fn := range fns {
		fn(zc)
	}

	// create a default context provider.
	if zc.ctx == nil {
		zc.ctx, zc.ctxCanceler = context.WithCancel(context.Background())
	}

	if zc.readRequests == nil {
		zc.readRequests = make(ZPayloadStream)
	}

	if zc.writeRequests == nil {
		zc.writeRequests = make(ZPayloadStream)
	}

	switch zc.conn.(type) {
	case *net.UDPConn:
		if zc.readWorker == nil {
			zc.readWorker = zc.readTCPFrom
		}

		if zc.writeWorker == nil {
			zc.writeWorker = zc.writeToTCP
		}
	case *net.TCPConn:
		if zc.readWorker == nil {
			zc.readWorker = zc.readUDPFrom
		}

		if zc.writeWorker == nil {
			zc.writeWorker = zc.writeToUDP
		}
	default:
		if zc.readWorker == nil {
			zc.readWorker = zc.readUDPFrom
		}

		if zc.writeWorker == nil {
			zc.writeWorker = zc.writeToUDP
		}
	}

	zc.streamReader = &nbytes.DelimitedStreamReader{
		Src:        zc.conn,
		ReadBuffer: int(zc.readBuffer),
		Escape:     []byte(defaultEscape),
		Delimiter:  []byte(defaultDelimiter),
	}

	zc.streamWriter = &nbytes.DelimitedStreamWriter{
		Dest:        zc.conn,
		WriteBuffer: int(zc.writeBuffer),
		Escape:      []byte(defaultEscape),
		Delimiter:   []byte(defaultDelimiter),
	}

	return zc
}

// Reads returns the underline channel used for receiving new incoming reads.
func (zc *ZConn) Reads() chan *ZPayload {
	return zc.readRequests
}

// Writes returns the underline channel used for receiving write requests.
func (zc *ZConn) Writes() chan *ZPayload {
	return zc.writeRequests
}

// RemoteAddr returns giving remote net.Addr of ZConn.
func (zc *ZConn) RemoteAddr() net.Addr {
	return zc.addr
}

// LocalAddr returns giving net.Addr of ZConn.
func (zc *ZConn) LocalAddr() net.Addr {
	return zc.laddr
}

func (zc *ZConn) Close() error {
	if zc.ctxCanceler != nil {
		zc.ctxCanceler()
	}
	return nil
}

func (zc *ZConn) writeLoop() {
	zc.waiter.Add(1)
	go func() {
		defer zc.waiter.Done()

		var done = zc.ctx.Done()

		for {
			select {
			case <-done:
				// we are being asked to stop and close.
				return
			case req := <-zc.writeRequests:
				if req == nil {
					select {
					case <-done:
						return
					default:
						continue
					}
				}

				if err := zc.writeWorker(zc.streamWriter, req); err != nil {
					log.Printf("[ZConn] | Failed connection writing process: %s", err)

					if err == ErrKillConnection {
						zc.ctxCanceler()
						ReleaseZPayload(req)
						return
					}

					select {
					case <-done:
						return
					default:
						continue
					}
				}
			}
		}
	}()
}

func (zc *ZConn) writeToUDP(src io.Writer, zp *ZPayload) error {

	return nil
}

func (zc *ZConn) readUDPFrom(src io.Reader, zp *ZPayload) error {

	return nil
}

// readLoop lunches underline read loop.
func (zc *ZConn) readLoop() {
	zc.waiter.Add(1)
	go func() {
		defer zc.waiter.Done()

		var done = zc.ctx.Done()

		for {
			select {
			case <-done:
				// we are being asked to stop and close.
				return
			case req := <-zc.readRequests:
				if req == nil {
					select {
					case <-done:
						return
					default:
						continue
					}
				}

				if err := zc.readWorker(zc.streamReader, req); err != nil {
					log.Printf("[ZConn] | Failed connection reading process: %s", err)

					if req.Err != nil {
						req.Err <- err
					}

					if err == ErrKillConnection {
						zc.ctxCanceler()
						return
					}

					select {
					case <-done:
						return
					default:
						continue
					}
				}
			}
		}
	}()
}

func (zc *ZConn) writeToTCP(dest io.Writer, zp *ZPayload) error {
	if _, err := io.Copy(dest, zp.Stream); err != nil {
		if zp.Err != nil {
			zp.Err <- err
		}
		return err
	}

	if err := zp.Stream.Close(); err != nil {
		if zp.Err != nil {
			zp.Err <- err
		}
		return err
	}

	if zp.Err != nil {
		close(zp.Err)
	}
	return nil
}

func (zc *ZConn) readTCPFrom(src io.Reader, zp *ZPayload) error {
	if zp.Addr != nil {
		zp.Addr <- zc.addr
	}

	var err error
	for {
		_, err = io.Copy(zp.Stream, src)
		if err != nil {
			if err == nbytes.ErrEOS {
				break
			}

			log.Printf("[ZConn] | Failed to finish readFor: %s", err)
			return err
		}
	}
	return nil
}

func sameTimeout(t time.Duration) ZTimeout {
	return func(_ time.Duration) time.Duration {
		return t
	}
}

func copyByteSlice(b []byte, factor int) []byte {
	var nextCap = gradualExpansion(cap(b), len(b), factor)
	var newList = make([]byte, nextCap)
	var copied = copy(newList, b)
	return newList[:copied]
}

func gradualExpansion(capacity int, last int, factor int) int {
	if last == 0 {
		last = 1
	}
	var pb = capacity * factor
	var inc = (pb / capacity) + factor
	return inc + (capacity / last)
}
