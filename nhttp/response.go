package nhttp

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

var (
	// ErrNoPush is returned when underline connection does
	// not support Push API.
	ErrNoPush = errors.New("push not supported")
)

// Response wraps an http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See: https://golang.org/pkg/net/http/#ResponseWriter
type Response struct {
	beforeFuncs []func()
	afterFuncs  []func()
	Writer      http.ResponseWriter
	Status      int
	Size        int64
	finished    bool
}

// Before adds the giving function into a response before list.
// Which should be executed before any write, but if the response
// was already completed then the function is called.
func (r *Response) Before(fn func()) {
	if r.finished {
		fn()
		return
	}
	r.beforeFuncs = append(r.beforeFuncs, fn)
}

// After adds the giving function into a response after list.
func (r *Response) After(fn func()) {
	if r.finished {
		fn()
		return
	}
	r.afterFuncs = append(r.afterFuncs, fn)
}

// Header returns the header map for the writer that will be sent by
// WriteHeader. Changing the header after a call to WriteHeader (or Write) has
// no effect unless the modified headers were declared as trailers by setting
// the "Trailer" header before the call to WriteHeader (see example)
// To suppress implicit response headers, set their value to nil.
// Example: https://golang.org/pkg/net/http/#example_ResponseWriter_trailers
func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(code int) {
	if r.finished {
		return
	}

	for _, fn := range r.beforeFuncs {
		fn()
	}

	r.Status = code
	r.Writer.WriteHeader(code)
	r.finished = true

	for _, fn := range r.afterFuncs {
		fn()
	}
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *Response) Write(b []byte) (n int, err error) {
	if !r.finished {
		r.WriteHeader(http.StatusOK)
	}

	n, err = r.Writer.Write(b)
	r.Size += int64(n)
	return
}

// Push adds support for http.Pusher, if available and lets you push resources.
func (r *Response) Push(target string, ops *http.PushOptions) error {
	if pusher, ok := r.Writer.(http.Pusher); ok {
		return pusher.Push(target, ops)
	}
	return ErrNoPush
}

// Sent returns true/false if giving response has being sent/written.
func (r *Response) Sent() bool {
	return r.finished
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.Writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.Writer.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface to allow detecting
// when the underlying connection has gone away.
// This mechanism can be used to cancel long operations on the server if the
// client has disconnected before the response is ready.
// See [http.CloseNotifier](https://golang.org/pkg/net/http/#CloseNotifier)
func (r *Response) CloseNotify() <-chan bool {
	return r.Writer.(http.CloseNotifier).CloseNotify()
}

func (r *Response) reset(w http.ResponseWriter) {
	r.Writer = w
	r.Size = 0
	r.Status = http.StatusOK
	r.finished = false
}
