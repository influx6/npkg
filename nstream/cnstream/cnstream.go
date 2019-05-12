package cnstream

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gokit/npkg/nerror"
	"github.com/gokit/npkg/nstream"
)

var (
	noTime = time.Time{}
)

//*********************************************************************************************
// TimedConn
//*********************************************************************************************

// TimedConn implements a wrapper around a net.Conn which guards giving connection
// with appropriate read/write timeout.
type TimedConn struct {
	conn         net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewTimedConn returns a new instance of a TimedConn.
func NewTimedConn(conn net.Conn, rd time.Duration, wd time.Duration) *TimedConn {
	return &TimedConn{
		conn:         conn,
		readTimeout:  rd,
		writeTimeout: wd,
	}
}

// Write calls the underline connection read with provided timeout.
func (c *TimedConn) Write(b []byte) (int, error) {
	var writeErr = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if writeErr != nil {
		return 0, writeErr
	}

	var writeCount, err = c.conn.Write(b)
	if err != nil {
		return writeCount, err
	}

	var resetErr = c.conn.SetWriteDeadline(noTime)
	if resetErr != nil {
		return writeCount, resetErr
	}

	return writeCount, nil
}

// Read calls the underline connection read with provided timeout.
func (c *TimedConn) Read(b []byte) (int, error) {
	var readErr = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if readErr != nil {
		return 0, readErr
	}

	var readCount, err = c.conn.Read(b)
	if err != nil {
		return readCount, err
	}

	var resetErr = c.conn.SetReadDeadline(noTime)
	if resetErr != nil {
		return readCount, resetErr
	}

	return readCount, nil
}

/**********************************************************************************************

MessageRetrieval Message: NX {MAX_64_INT}

**********************************************************************************************/

type ReadConnection struct{
	addr net.Addr
	conn net.Conn
}

//*********************************************************************************************
// ReadSubscription
//*********************************************************************************************

// ReadSubscription (ConnectionSubscription) implements a subscription wrapper over a provided net.Conn.
type ReadSubscription struct {
	err       error
	bp        uint
	requested uint64
	ctr       sync.Cond
	sub       nstream.ByteSubscriber
	cl        sync.Mutex
	buffer    *bufio.Reader
	conn      *TimedConn
}

// NewReadSubscription returns a new ReadSubscription for giving connection and subscriber.
func NewReadSubscription(buffer uint, conn *TimedConn, subscriber nstream.ByteSubscriber) *ReadSubscription {
	var cn ReadSubscription
	cn.buffer = bufio.NewReaderSize(conn, int(buffer))
	cn.sub = subscriber
	cn.conn = conn
	cn.bp = buffer
	cn.manage()
	return &cn
}

// Next requests new the next amount of bytes to be retrieved from connection if
// available.
func (c *ReadSubscription) Next(n uint64) error {
	if c.err != nil {
		return c.err
	}

	if !c.isValid() {
		return nerror.New("subscription is inactive")
	}

	if n == 0 {
		return nil
	}

	atomic.AddUint64(&c.requested, n)
	c.ctr.Broadcast()
	return nil
}

// isValid returns true whether giving subscription is active
// or not.
func (c *ReadSubscription) isValid() bool {
	c.cl.Lock()
	var hasConn = c.conn == nil
	c.cl.Unlock()
	return hasConn
}

func (c *ReadSubscription) manage() {
	var buffer []byte
	go func() {
		for {
			c.ctr.L.Lock()
			c.ctr.Wait()
			c.ctr.L.Unlock()

			var requested = uint(atomic.LoadUint64(&c.requested))
			if requested == 0 {
				continue
			}
			if requested > c.bp {
				requested = c.bp
			}

			buffer = make([]byte, int(requested))
			read, err := c.conn.Read(buffer)
			if err != nil {
				if tmpErr, ok := err.(net.Error); ok {
					if tmpErr.Timeout() || tmpErr.Temporary() {
						continue
					}
				}
			}

			if read <= 0 {
				continue
			}

			buffer = buffer[:read]
			if uint(read) == requested {
				atomic.StoreUint64(&c.requested, 0)
				continue
			}

			var left = uint64(requested - uint(read))
			atomic.AddUint64(&c.requested, left)
			c.sub.OnNext(buffer)
		}
	}()
}
