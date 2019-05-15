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

	"golang.org/x/sync/errgroup"

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
// Server
//**********************************************************************

// ServerHandler defines an interface type for handling connection
// in isolated goroutines.
type ServerHandler interface {
	ServeConn(context.Context, net.Conn) error
}

// Server implements a basic, structured wrapper around a net.Listener
// to handle and process incoming connections.
type Server struct {
	Handler  ServerHandler
	Listener net.Listener
	ctx      context.Context
	waiter   *errgroup.Group
}

// NewServer returns a new instance of a Server.
func NewServer(ctx context.Context, handler ServerHandler, listener net.Listener) *Server {
	var gp *errgroup.Group
	gp, ctx = errgroup.WithContext(ctx)

	return &Server{
		ctx:      ctx,
		waiter:   gp,
		Handler:  handler,
		Listener: listener,
	}
}

// Serve initializes the connection handling and reading goroutine.
// It returns immediately without block. It depends on you to use the
// Server.Wait to block as needed till server is closed.
func (s *Server) Serve() {
	s.waiter.Go(s.serviceListener)
}

// Wait blocks underline server till connection
// has returned. It returns any possible error that
// occurred.
func (s *Server) Wait() error {
	return s.waiter.Wait()
}

// serviceListener blocks and handles incoming connections
// by server. It will only ever stop working if accepts
// failed or the Handler.ServerConn returns a ErrKillConnection
// error.
func (s *Server) serviceListener() error {
	for {
		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		default:
		}

		conn, err := s.Listener.Accept()
		if err != nil {
			if tmpErr, ok := err.(net.Error); ok && tmpErr.Temporary() {
				continue
			}
			return err
		}

		if err := s.Handler.ServeConn(s.ctx, conn); err != nil {
			if err == ErrKillConnection {
				return err
			}

			continue
		}
	}
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
			payload.Addr = make(chan net.Addr, 1)
			payload.Done = make(chan struct{}, 1)
			return payload
		},
	}
)

// AcquireZPayload returns a new ZPayload for use.
func AcquireZPayload() *ZPayload {
	if nzp, ok := zPayloads.Get().(*ZPayload); ok {
		nzp.Reset()
		return nzp
	}

	var payload ZPayload
	payload.Err = make(chan error, 1)
	payload.Addr = make(chan net.Addr, 1)
	payload.Done = make(chan struct{}, 1)
	return &payload
}

// ReleaseZPayload returns giving ZPayload into underline object pool
// for reuse.
// Never re-use a ZPayload reference you've previously released.
func ReleaseZPayload(zp *ZPayload) {
	zp.Clear()
	zPayloads.Put(zp)
}

// ZConnWorker defines a interface type for the servicing of
// an underline read and write requests.
type ZConnWorker interface {
	ServeRead(context.Context, io.Reader, *ZPayload) error
	ServeWrite(context.Context, io.Writer, *ZPayload) error
}

// ZPayload defines an underline structure for writing data
// into an underline ZConn.
type ZPayload struct {
	Done   chan struct{}
	Err    chan error
	Addr   chan net.Addr
	Stream io.ReadWriteCloser
}

// Clear clears ZPayload fields.
func (z *ZPayload) Clear() {
	z.Err = nil
	z.Done = nil
	z.Addr = nil
	z.Stream = nil
}

// Reset giving ZPayload fields.
func (z *ZPayload) Reset() {
	z.Err = make(chan error, 1)
	z.Addr = make(chan net.Addr, 1)
	z.Done = make(chan struct{}, 1)
}

func (z *ZPayload) verify() {
	if z.Stream == nil {
		panic("ZPayload.Stream can not be nil")
	}
	if z.Done == nil {
		panic("ZPayload.Done can not be nil")
	}
	if z.Err == nil {
		panic("ZPayload.Err can not be nil")
	}
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
	worker              ZConnWorker
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
func NewZConn(conn net.Conn, worker ZConnWorker, fns ...ZApply) *ZConn {
	var zc = new(ZConn)
	zc.conn = conn
	zc.worker = worker
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

	if zc.ctx == nil {
		zc.ctx, zc.ctxCanceler = context.WithCancel(context.Background())
	}

	if zc.readRequests == nil {
		zc.readRequests = make(ZPayloadStream)
	}

	if zc.writeRequests == nil {
		zc.writeRequests = make(ZPayloadStream)
	}

	switch conn.(type) {
	case *net.TCPConn:
		zc.worker = &TCPWorker{}
	case *net.UDPConn:
		zc.worker = &UDPWorker{}
	default:
		panic("ZConn.Worker must be provided as it's a hard requirement")
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

	// boot up read loop.
	zc.readLoop()

	// boot up write loop.
	zc.writeLoop()

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

	zc.waiter.Wait()
	return nil
}

type TCPWorker struct{}

// ServeRead handles servicing a read request against provided io.Reader which is
// read into the underline connection.
func (zc TCPWorker) ServeRead(ctx context.Context, src io.Reader, zp *ZPayload) error {
	fmt.Printf("Will Read \n\n")

	var read, err = io.Copy(zp.Stream, src)
	fmt.Printf("Read : %#v -> %s\n\n\n", read, err)
	if err != nil {
		log.Printf("[ZConn] | Failed to finish readFor: %s", err)

		if err == nbytes.ErrEOS {
			log.Printf("[ZConn] | Read %d bytes from connection", read)
			return nil
		}

		return err
	}

	log.Printf("[ZConn.readTCPFrom] | Read %d bytes from connection", read)
	return nil
}

// ServeRead handles servicing a read request against provided io.Reader which is
// read into the underline connection.
func (zc TCPWorker) ServeWrite(ctx context.Context, src io.Writer, zp *ZPayload) error {
	log.Printf("[ZConn] | Written new stream into connection")

	var written, err = io.Copy(dest, zp.Stream)
	if err != nil {
		return err
	}

	log.Printf("[ZConn] | Written %d bytes into connection", written)

	if err := zp.Stream.Close(); err != nil {
		return err
	}

	return nil
}

type UDPWorker struct{}

// ServeRead handles servicing a read request against provided io.Reader which is
// read into the underline connection.
func (zc UDPWorker) ServeRead(ctx context.Context, src io.Reader, zp *ZPayload) error {

	return nil
}

// ServeRead handles servicing a read request against provided io.Reader which is
// read into the underline connection.
func (zc UDPWorker) ServeWrite(ctx context.Context, src io.Writer, zp *ZPayload) error {

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
			case req, ok := <-zc.writeRequests:
				if !ok {
					return
				}

				// verify giving request object is valid.
				req.verify()

				if err := zc.worker.ServeWrite(zc.ctx, zc.streamWriter, req); err != nil {
					log.Printf("[ZConn] | Failed connection writing process: %s", err)

					if req.Err != nil {
						req.Err <- err
					}

					if err == ErrKillConnection {
						zc.ctxCanceler()
						close(req.Done)
						return
					}

					continue
				}

				close(req.Done)
			}
		}
	}()
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
			case req, ok := <-zc.readRequests:
				if !ok {
					return
				}

				// verify giving request object is valid.
				req.verify()

				if req.Addr != nil {
					req.Addr <- zc.addr
				}

				if err := zc.worker.ServeRead(zc.ctx, zc.streamReader, req); err != nil {
					log.Printf("[ZConn] | Failed connection reading process: %s", err)

					if req.Err != nil {
						req.Err <- err
					}

					if err == ErrKillConnection {
						zc.ctxCanceler()
						close(req.Done)
						return
					}

					continue
				}

				close(req.Done)
			}
		}
	}()
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
